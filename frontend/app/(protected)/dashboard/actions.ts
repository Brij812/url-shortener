"use server";

import axios from "axios";
import { cookies } from "next/headers";

// 🟢 Fetch all URLs
export async function getUserUrlsAction() {
  const api = process.env.NEXT_PUBLIC_API_BASE_URL!;
  const cookieStore = await cookies();
  const token = cookieStore.get("hl_jwt")?.value;

  try {
    const res = await axios.get(`${api}/all`, {
      headers: { Authorization: `Bearer ${token}` },
      validateStatus: () => true,
    });

    if (res.status !== 200) {
      return { error: res.data?.message || "Failed to fetch URLs" };
    }

    return { urls: res.data };
  } catch (err) {
    console.error("❌ getUserUrlsAction error:", err);
    return { error: "Something went wrong while fetching URLs" };
  }
}

// 🔴 Delete a short URL
export async function deleteUrlAction(code: string) {
  const api = process.env.NEXT_PUBLIC_API_BASE_URL!;
  const cookieStore = await cookies();
  const token = cookieStore.get("hl_jwt")?.value;

  try {
    console.log(`🗑️ Deleting short URL: ${code}`);

    const res = await axios.delete(`${api}/url/${code}`, {
      headers: { Authorization: `Bearer ${token}` },
      validateStatus: () => true,
    });

    console.log("🗑️ Delete response:", res.status, res.data);

    if (res.status !== 200) {
      return { error: res.data?.message || "Failed to delete link" };
    }

    return { success: true };
  } catch (err) {
    console.error("❌ deleteUrlAction error:", err);
    return { error: "Something went wrong while deleting link" };
  }
}
