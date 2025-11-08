"use server";

import axios from "axios";
import { cookies } from "next/headers";

export async function getMetricsAction() {
  const api = process.env.NEXT_PUBLIC_API_BASE_URL!;
  const cookieStore = await cookies();
  const token = cookieStore.get("hl_jwt")?.value;

  try {
    // 1️⃣ Get metrics (domain → count)
    const metricsRes = await axios.get(`${api}/metrics`, {
      headers: { Authorization: `Bearer ${token}` },
      validateStatus: () => true,
    });

    // 2️⃣ Get all shortened URLs to count total unique links
    const allRes = await axios.get(`${api}/all`, {
      headers: { Authorization: `Bearer ${token}` },
      validateStatus: () => true,
    });

    console.log("📊 Raw backend metrics response:", metricsRes.data);
    console.log("🔗 All shortened URLs response:", allRes.data);

    if (metricsRes.status !== 200) {
      return { error: metricsRes.data?.message || "Failed to load metrics" };
    }

    let cleanData: any[] = [];

    // 🧩 Convert object {domain: count, ...} → array [{domain, count}]
    if (Array.isArray(metricsRes.data)) {
      cleanData = metricsRes.data;
    } else if (typeof metricsRes.data === "object" && metricsRes.data !== null) {
      cleanData = Object.entries(metricsRes.data).map(([domain, count]) => ({
        domain,
        count: Number(count),
      }));
    }

    const totalShortened = Array.isArray(allRes.data) ? allRes.data.length : 0;

    return { data: cleanData, totalShortened };
  } catch (err) {
    console.error("❌ getMetricsAction error:", err);
    return { error: "Something went wrong" };
  }
}
