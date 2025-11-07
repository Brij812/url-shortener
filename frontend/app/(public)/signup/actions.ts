"use server";

import axios from "axios";

export async function signupAction(data: { email: string; password: string }) {
  console.log("API URL =", process.env.NEXT_PUBLIC_API_BASE_URL);

  try {
    const res = await axios.post(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/signup`,
      data,
      { validateStatus: () => true }
    );

    console.log("Signup response:", res.status, res.data);

    // ✅ Accept all success 2xx codes
    if (res.status < 200 || res.status >= 300) {
      return { error: res.data?.message || "Signup failed" };
    }

    return { success: true };
  } catch (err) {
    console.error("Signup error:", err);
    return { error: "Something went wrong" };
  }
}
