"use client";

import { useEffect, useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { toast } from "sonner";
import { getUserUrlsAction, deleteUrlAction } from "./actions";
import { Copy, ExternalLink, Trash2, X } from "lucide-react";

export default function DashboardPage() {
  const [urls, setUrls] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedCode, setSelectedCode] = useState<string | null>(null);

  useEffect(() => {
    async function fetchData() {
      try {
        const res = await getUserUrlsAction();
        if (res.error) toast.error(res.error);
        else setUrls(res.urls || []);
      } catch (err) {
        console.error("❌ Fetch error:", err);
        toast.error("Failed to load URLs");
      } finally {
        setLoading(false);
      }
    }
    fetchData();
  }, []);

  async function handleDelete() {
    if (!selectedCode) return;
    const res = await deleteUrlAction(selectedCode);
    if (res.error) return toast.error(res.error);

    toast.success("Link deleted!");
    setUrls((prev) =>
      prev.filter((u) => {
        const c = u.code || u.short_code || u.short_url?.split("/").pop();
        return c !== selectedCode;
      })
    );
    setSelectedCode(null);
  }

  return (
    <div className="p-10 relative">
      <motion.h1
        initial={{ opacity: 0, y: -10 }}
        animate={{ opacity: 1, y: 0 }}
        className="text-3xl font-black mb-6"
      >
        Your Shortened URLs
      </motion.h1>

      {loading ? (
        <p className="text-gray-500 animate-pulse">Loading URLs...</p>
      ) : urls.length === 0 ? (
        <p className="text-gray-500">You have not created any short URLs yet.</p>
      ) : (
        <motion.table
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          className="w-full border-collapse border border-gray-300 shadow-sm rounded-xl overflow-hidden"
        >
          <thead className="bg-gray-100">
            <tr>
              <th className="border border-gray-300 px-4 py-2 text-left font-semibold w-1/6">
                Short Code
              </th>
              <th className="border border-gray-300 px-4 py-2 text-left font-semibold w-3/6">
                Long URL
              </th>
              <th className="border border-gray-300 px-4 py-2 text-left font-semibold w-1/6">
                Created
              </th>
              <th className="border border-gray-300 px-4 py-2 text-left font-semibold w-1/6">
                Actions
              </th>
            </tr>
          </thead>
          <tbody>
            {urls.map((u, i) => {
              const code =
                u.code || u.short_code || u.short_url?.split("/").pop() || "—";
              const longUrl = u.url || u.long_url || "—";
              const created = u.created_at
                ? new Date(u.created_at).toLocaleString()
                : "—";
              const fullShortUrl =
                u.short_url || `http://localhost:8080/${code}`;

              return (
                <tr
                  key={i}
                  className="border border-gray-300 hover:bg-gray-50 transition"
                >
                  <td className="p-3 font-mono text-blue-700">{code}</td>
                  <td className="p-3 break-all text-blue-600">
                    <a
                      href={longUrl}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="hover:underline"
                    >
                      {longUrl}
                    </a>
                  </td>
                  <td className="p-3 text-gray-700">{created}</td>
                  <td className="p-3">
                    <div className="flex items-center gap-4 justify-center">
                      <button
                        onClick={() => {
                          navigator.clipboard.writeText(fullShortUrl);
                          toast.success("Copied!");
                        }}
                        title="Copy Short URL"
                        className="p-2 rounded hover:bg-gray-200 transition"
                      >
                        <Copy size={18} />
                      </button>

                      <a
                        href={fullShortUrl}
                        target="_blank"
                        className="p-2 rounded hover:bg-gray-200 transition"
                        title="Open Short URL"
                      >
                        <ExternalLink size={18} />
                      </a>

                      <button
                        onClick={() => setSelectedCode(code)}
                        title="Delete URL"
                        className="p-2 rounded hover:bg-red-100 transition"
                      >
                        <Trash2
                          size={18}
                          className="text-red-600 hover:text-red-800"
                        />
                      </button>
                    </div>
                  </td>
                </tr>
              );
            })}
          </tbody>
        </motion.table>
      )}

      {/* 🧩 Delete Confirmation Modal */}
      <AnimatePresence>
        {selectedCode && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/40 backdrop-blur-sm flex items-center justify-center z-50"
          >
            <motion.div
              initial={{ scale: 0.9, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.9, opacity: 0 }}
              className="bg-white border-2 border-black rounded-xl shadow-lg max-w-sm w-full p-6 text-center relative"
            >
              <button
                className="absolute top-3 right-3 text-gray-500 hover:text-black"
                onClick={() => setSelectedCode(null)}
              >
                <X size={20} />
              </button>

              <h2 className="text-xl font-bold mb-2">Delete Link?</h2>
              <p className="text-gray-600 mb-6">
                Are you sure you want to delete this short link?
              </p>

              <div className="flex justify-center gap-4">
                <button
                  onClick={() => setSelectedCode(null)}
                  className="px-4 py-2 bg-gray-300 hover:bg-gray-400 rounded font-semibold transition"
                >
                  Cancel
                </button>
                <button
                  onClick={handleDelete}
                  className="px-4 py-2 bg-red-600 text-white hover:bg-red-700 rounded font-semibold transition"
                >
                  Delete
                </button>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
