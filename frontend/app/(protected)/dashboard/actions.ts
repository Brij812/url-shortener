"use server";

import axios from "axios";
import { cookies } from "next/headers";

export async function getUserUrlsAction() {
  const api = process.env.NEXT_PUBLIC_API_BASE_URL!;
  const cookieStore = await cookies();
  const token = cookieStore.get("hl_jwt")?.value;

  try {
    const res = await axios.get(`${api}/all`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
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
