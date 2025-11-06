"use server";

import { cookies } from "next/headers";
import axios from "axios";

export async function loginAction(data: { email: string; password: string }) {
  try {
    const res = await axios.post(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/login`,
      data,
      {
        validateStatus: () => true,
      }
    );

    if (res.status !== 200) {
      return { error: res.data?.message || "Invalid email/password" };
    }

    const token = res.data.token;

    if (!token) {
      return { error: "Token missing in response" };
    }

    const cookieStore = await cookies();

    cookieStore.set({
      name: process.env.COOKIE_NAME || "hl_jwt",
      value: token,
      httpOnly: true,
      path: "/",
      sameSite: "lax",
      secure: process.env.NODE_ENV === "production",
      maxAge: 60 * 60, // 1 hour
    });

    return { success: true };
  } catch (err) {
    console.error("Login error:", err);
    return { error: "Something went wrong" };
  }
}
