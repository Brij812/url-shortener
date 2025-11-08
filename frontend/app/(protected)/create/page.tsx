"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { createUrlAction } from "./actions";
import { toast } from "sonner";
import { motion } from "framer-motion";
import { Copy, ExternalLink, Loader2 } from "lucide-react";

const formSchema = z.object({
  url: z.string().url("Enter a valid URL"),
});

type FormValues = z.infer<typeof formSchema>;

export default function CreatePage() {
  const [result, setResult] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: { url: "" },
  });

  async function onSubmit(values: FormValues) {
    try {
      setLoading(true);
      setResult(null);

      console.log("➡️ Sending create request:", values);

      const raw = await createUrlAction(values);
      const res = typeof raw === "string" ? JSON.parse(raw) : raw;

      if (res.error) {
        toast.error(res.error);
        return;
      }

      const shortUrl = res.short_url || res.short_code;
      if (!shortUrl) {
        toast.error("Failed to get short URL");
        return;
      }

      setResult(shortUrl);
      toast.success("Short URL created!");
      form.reset();
    } catch (err) {
      console.error("❌ Create error:", err);
      toast.error("Something went wrong.");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 flex flex-col items-center justify-center px-6">
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        className="w-full max-w-xl bg-white rounded-2xl border-2 border-gray-800 shadow-xl p-8"
      >
        <motion.h1
          initial={{ opacity: 0, y: -10 }}
          animate={{ opacity: 1, y: 0 }}
          className="text-4xl font-black mb-8 text-center text-gray-900"
        >
          Create Short URL
        </motion.h1>

        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
          <div className="flex flex-col space-y-2">
            <label className="text-lg font-semibold text-gray-800">Long URL</label>
            <input
              {...form.register("url")}
              placeholder="https://example.com/my-long-url"
              className="border-2 border-gray-800 p-3 text-lg font-mono rounded-md focus:outline-none focus:ring-2 focus:ring-gray-900 bg-gray-50"
            />
            {form.formState.errors.url && (
              <p className="text-red-600 text-sm">
                {form.formState.errors.url.message}
              </p>
            )}
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full bg-gray-900 text-white py-3 text-lg font-bold rounded-md hover:bg-gray-800 transition active:scale-95 flex items-center justify-center gap-2"
          >
            {loading && <Loader2 size={20} className="animate-spin" />}
            {loading ? "Creating..." : "Create"}
          </button>
        </form>

        {result && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            className="mt-8 border-t-2 border-gray-800 pt-4 space-y-3"
          >
            <p className="font-semibold text-lg text-gray-800 text-center">
              Your Short URL
            </p>

            <div className="flex items-center justify-between border-2 border-gray-800 p-3 font-mono break-all bg-gray-50 rounded-md">
              <span className="truncate text-gray-700">{result}</span>

              <div className="flex gap-3">
                <button
                  onClick={() => {
                    navigator.clipboard.writeText(result);
                    toast.success("Copied!");
                  }}
                  className="hover:text-blue-600 transition"
                >
                  <Copy size={20} />
                </button>

                <a
                  href={result}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="hover:text-blue-600 transition"
                >
                  <ExternalLink size={20} />
                </a>
              </div>
            </div>
          </motion.div>
        )}
      </motion.div>
    </div>
  );
}
