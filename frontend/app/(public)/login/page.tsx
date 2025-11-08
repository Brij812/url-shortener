"use client";

import { useState } from "react";
import { loginAction } from "./actions";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";

export default function LoginPage() {
  const [error, setError] = useState("");

  async function handleSubmit(formData: FormData) {
    setError("");
    const email = formData.get("email")?.toString();
    const password = formData.get("password")?.toString();

    if (!email || !password) {
      setError("Please fill in both fields.");
      return;
    }

    const res = await loginAction({ email, password });
    if (res?.error) {
      setError(res.error);
      toast.error(res.error);
    } else {
      toast.success("Welcome back!");
      window.location.href = "/dashboard";
    }
  }

  return (
    <div className="min-h-screen flex">
      {/* Left Section */}
      <div className="flex-1 bg-[#111827] text-white flex flex-col justify-center items-center p-12">
        <h1 className="text-4xl font-bold mb-4">Welcome Back 👋</h1>
        <p className="text-gray-300 max-w-md text-center">
          Log in to manage your links, track analytics, and grow your reach —
          all from one clean dashboard.
        </p>
      </div>

      {/* Right Section */}
      <div className="flex-1 flex justify-center items-center bg-white p-12">
        <form
          action={handleSubmit}
          className="w-full max-w-sm space-y-6 border border-gray-200 rounded-lg shadow-lg p-8"
        >
          <h2 className="text-2xl font-semibold text-center">Login</h2>

          <div>
            <label className="block text-sm font-medium text-gray-700">
              Email
            </label>
            <Input
              name="email"
              type="email"
              placeholder="you@example.com"
              className="mt-1"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700">
              Password
            </label>
            <Input
              name="password"
              type="password"
              placeholder="••••••••"
              className="mt-1"
              required
            />
          </div>

          {error && <p className="text-red-500 text-sm">{error}</p>}

          <Button className="w-full" type="submit">
            Login
          </Button>

          <p className="text-center text-sm text-gray-500 mt-4">
            Don’t have an account?{" "}
            <a href="/signup" className="text-blue-600 hover:underline">
              Sign up
            </a>
          </p>
        </form>
      </div>
    </div>
  );
}
