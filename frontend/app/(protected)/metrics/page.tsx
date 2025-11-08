"use client";

import { useEffect, useState } from "react";
import { getMetricsAction } from "./actions";
import { toast } from "sonner";
import { motion } from "framer-motion";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell,
} from "recharts";
import { Loader2 } from "lucide-react";

export default function MetricsPage() {
  const [data, setData] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [totalShortened, setTotalShortened] = useState(0);

  useEffect(() => {
    (async () => {
      const res = await getMetricsAction();
      if (res.error) {
        toast.error(res.error);
        setLoading(false);
        return;
      }

      const chartData = res.data || [];
      console.log("📈 Final metrics data for chart:", chartData);
      setData(chartData);
      setTotalShortened(res.totalShortened || 0);
      setLoading(false);
    })();
  }, []);

  const totalClicks = data.reduce((sum, d) => sum + d.count, 0);
  const topDomain = data.length ? data[0].domain : "-";

  const COLORS = ["#000000", "#444444", "#777777", "#aaaaaa", "#dddddd"];

  return (
    <div className="max-w-5xl mx-auto mt-12 p-4">
      <motion.h1
        initial={{ opacity: 0, y: -10 }}
        animate={{ opacity: 1, y: 0 }}
        className="text-4xl font-black mb-10 tracking-tight"
      >
        Metrics Overview
      </motion.h1>

      {loading ? (
        <div className="flex justify-center items-center h-80">
          <Loader2 size={32} className="animate-spin text-black" />
        </div>
      ) : data.length === 0 ? (
        <div className="text-center text-gray-500 text-lg mt-20">
          No activity yet 😴
        </div>
      ) : (
        <>
          {/* === Stats Cards === */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className="grid grid-cols-1 sm:grid-cols-3 gap-6 mb-10"
          >
            <div className="border-2 border-black p-6 text-center bg-white">
              <p className="text-sm font-medium uppercase text-gray-600">
                Total Clicks
              </p>
              <p className="text-3xl font-black mt-2">{totalClicks}</p>
            </div>

            <div className="border-2 border-black p-6 text-center bg-white">
              <p className="text-sm font-medium uppercase text-gray-600">
                Top Domain
              </p>
              <p className="text-2xl font-black mt-2 break-words">{topDomain}</p>
            </div>

            <div className="border-2 border-black p-6 text-center bg-white">
              <p className="text-sm font-medium uppercase text-gray-600">
                Total Shortened URLs
              </p>
              <p className="text-3xl font-black mt-2">{totalShortened}</p>
            </div>
          </motion.div>

          {/* === Charts === */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-10">
            {/* Bar Chart */}
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              className="border-2 border-black p-4 bg-white"
            >
              <h2 className="text-xl font-bold mb-4">Top Domains by Clicks</h2>
              <div className="h-80">
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart
                    data={data}
                    layout="vertical"
                    margin={{ top: 10, right: 30, left: 60, bottom: 10 }}
                  >
                    <XAxis type="number" />
                    <YAxis dataKey="domain" type="category" width={120} />
                    <Tooltip />
                    <Bar dataKey="count" fill="#000000" />
                  </BarChart>
                </ResponsiveContainer>
              </div>
            </motion.div>

            {/* Pie Chart */}
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              className="border-2 border-black p-4 bg-white"
            >
              <h2 className="text-xl font-bold mb-4">Domain Distribution</h2>
              <div className="h-80">
                <ResponsiveContainer width="100%" height="100%">
                  <PieChart>
                    <Pie
                      data={data}
                      dataKey="count"
                      nameKey="domain"
                      cx="50%"
                      cy="50%"
                      outerRadius={100}
                      label
                    >
                      {data.map((_, i) => (
                        <Cell key={i} fill={COLORS[i % COLORS.length]} />
                      ))}
                    </Pie>
                    <Tooltip />
                  </PieChart>
                </ResponsiveContainer>
              </div>
            </motion.div>
          </div>
        </>
      )}
    </div>
  );
}
