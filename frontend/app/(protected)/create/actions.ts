"use server";

import axios from "axios";
import { cookies } from "next/headers";

export async function createUrlAction(data: { url: string }) {
  const api = process.env.NEXT_PUBLIC_API_BASE_URL!;
  const cookieStore = await cookies();
  const token = cookieStore.get("hl_jwt")?.value;

  try {
    console.log("📡 Sending to backend:", `${api}/shorten`, "with", data);

    const res = await axios.post(`${api}/shorten`, JSON.stringify(data), {
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
      },
      validateStatus: () => true,
    });

    console.log("📦 Raw backend response data:", res.data);

    if (res.status !== 200 && res.status !== 201) {
      console.error("❌ Backend returned error:", res.status, res.data);
      return JSON.stringify({ error: res.data?.message || "Failed to create short URL" });
    }

    const short_url =
      res.data.short_url ||
      res.data.ShortURL ||
      res.data.url ||
      null;

    // 👇 Return as plain string (always serializable)
    return JSON.stringify({ short_url });
  } catch (err) {
    console.error("❌ createUrlAction error:", err);
    return JSON.stringify({ error: "Something went wrong" });
  }
}
