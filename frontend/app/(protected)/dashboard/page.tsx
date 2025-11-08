"use client";

import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { toast } from "sonner";
import { getUserUrlsAction } from "./actions";

export default function DashboardPage() {
  const [urls, setUrls] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

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

  return (
    <div className="p-10">
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
              <th className="border border-gray-300 px-4 py-2 text-left font-semibold">
                Short Code
              </th>
              <th className="border border-gray-300 px-4 py-2 text-left font-semibold">
                Long URL
              </th>
              <th className="border border-gray-300 px-4 py-2 text-left font-semibold">
                Created
              </th>
              <th className="border border-gray-300 px-4 py-2 text-left font-semibold">
                Actions
              </th>
            </tr>
          </thead>
          <tbody>
            {urls.map((u, i) => {
              // Normalize data shape
              const code = u.code || u.short_code || u.short_url?.split("/").pop() || "—";
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
                    <a
                      href={fullShortUrl}
                      target="_blank"
                      className="text-blue-600 hover:underline"
                    >
                      Open
                    </a>
                  </td>
                </tr>
              );
            })}
          </tbody>
        </motion.table>
      )}
    </div>
  );
}
