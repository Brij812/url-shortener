"use server";

import axios from "axios";

export async function signupAction(data: {
  email: string;
  password: string;
}) {
  try {
    const res = await axios.post(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/signup`,
      data,
      {
        validateStatus: () => true,
      }
    );

    if (res.status !== 200) {
      return { error: res.data?.message || "Signup failed" };
    }

    return { success: true };
  } catch (err) {
    console.error("Signup error:", err);
    return { error: "Something went wrong" };
  }
}
