"use client";

import { useEffect, useMemo, useState } from "react";
import axios from "axios";
import {AnimatePresence, motion } from "framer-motion";
import { toast } from "sonner";
import { LogOut, Moon, Bell, Settings } from "lucide-react";
import { X } from "lucide-react";


const API_BASE = process.env.NEXT_PUBLIC_API_BASE_URL!;

function getCookie(name: string): string | null {
  if (typeof document === "undefined") return null;
  const m = document.cookie.match(new RegExp(`(?:^|; )${name}=([^;]*)`));
  return m ? decodeURIComponent(m[1]) : null;
}

export default function SettingsPage() {
  const email = useMemo(() => getCookie("hl_email") || "user@hyperlinkos.dev", []);
  const [darkMode, setDarkMode] = useState(false);
  const [notifications, setNotifications] = useState(false);

  // 🌙 Theme setup
  useEffect(() => {
    const saved = typeof window !== "undefined" ? localStorage.getItem("hl_theme") : null;
    const isDark = saved === "dark";
    setDarkMode(isDark);
    document.documentElement.classList.toggle("dark", isDark);
  }, []);

  const handleDarkToggle = () => {
    const next = !darkMode;
    setDarkMode(next);
    document.documentElement.classList.toggle("dark", next);
    localStorage.setItem("hl_theme", next ? "dark" : "light");
    toast.success(`Dark mode ${next ? "enabled" : "disabled"}`);
  };

  const handleNotifToggle = () => {
    const next = !notifications;
    setNotifications(next);
    toast.success(
      next ? `Notifications enabled for ${email}` : "Notifications disabled"
    );
  };

  const handleChangePassword = () =>
    toast.info("Password reset flow coming soon.");

  // --- Add this near top of component ---
    const [confirmLogout, setConfirmLogout] = useState(false);

// --- Update your logout function ---
    const handleLogout = async () => {
     try {
    await axios.post(`${API_BASE}/logout`, {}, { withCredentials: true });
    localStorage.removeItem("hl_theme");
    document.cookie = "hl_email=; Max-Age=0; path=/";
    toast.success("Logged out successfully!");
    setTimeout(() => (window.location.href = "/login"), 600);
     } catch (err) {
    console.error("Logout failed:", err);
    toast.error("Failed to logout.");
     }
    };


  return (
    <div className="flex justify-center min-h-[90vh] bg-gray-50 dark:bg-neutral-950 p-6">
      <motion.div
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.35 }}
        className="w-full max-w-3xl bg-white dark:bg-neutral-900 rounded-2xl shadow-md border border-gray-200 dark:border-neutral-800 p-8 space-y-8"
      >
        <div className="text-center">
          <h1 className="text-4xl font-extrabold mb-1 text-neutral-900 dark:text-neutral-100">
            Settings
          </h1>
          <p className="text-gray-500 dark:text-neutral-400 text-sm">
            Manage your account and preferences.
          </p>
        </div>

        <div className="flex items-center gap-5 p-5 border border-gray-200 dark:border-neutral-800 rounded-xl">
          <div className="h-14 w-14 flex items-center justify-center bg-gradient-to-tr from-black to-gray-600 text-white font-bold text-xl rounded-full">
            {email.charAt(0).toUpperCase()}
          </div>
          <div>
            <h2 className="font-semibold text-lg">{email}</h2>
            <p className="text-gray-500 dark:text-neutral-400 text-sm">Signed in</p>
          </div>
        </div>

        <div className="grid md:grid-cols-2 gap-6">
          {/* Account */}
          <div className="border border-gray-200 dark:border-neutral-800 rounded-xl p-5 space-y-3">
            <h3 className="font-semibold text-lg flex items-center gap-2">
              <Settings size={18} /> Account
            </h3>
            <p className="text-sm text-gray-500">Manage your login credentials and profile info.</p>
            <button
              onClick={handleChangePassword}
              className="w-full mt-3 border border-black dark:border-white text-black dark:text-white py-2 rounded-md font-medium hover:bg-black hover:text-white dark:hover:bg-white dark:hover:text-black transition"
            >
              Change Password
            </button>
          </div>

          {/* Preferences */}
          <div className="border border-gray-200 dark:border-neutral-800 rounded-xl p-5 space-y-4">
            <h3 className="font-semibold text-lg flex items-center gap-2">
              <Moon size={18} /> Preferences
            </h3>

            <div className="flex justify-between">
              <span className="text-sm font-medium">Dark Mode</span>
              <button
                onClick={handleDarkToggle}
                className={`relative inline-flex h-6 w-11 items-center rounded-full transition ${
                  darkMode ? "bg-black" : "bg-gray-300"
                }`}
              >
                <span
                  className={`inline-block h-5 w-5 transform rounded-full bg-white transition ${
                    darkMode ? "translate-x-5" : "translate-x-1"
                  }`}
                />
              </button>
            </div>

            <div className="flex justify-between">
              <span className="text-sm font-medium">Notifications</span>
              <button
                onClick={handleNotifToggle}
                className={`relative inline-flex h-6 w-11 items-center rounded-full transition ${
                  notifications ? "bg-black" : "bg-gray-300"
                }`}
              >
                <span
                  className={`inline-block h-5 w-5 transform rounded-full bg-white transition ${
                    notifications ? "translate-x-5" : "translate-x-1"
                  }`}
                />
              </button>
            </div>
          </div>
        </div>

        {/* Logout Section */}
<div className="border-t dark:border-neutral-800 pt-6 text-center relative">
  <button
    onClick={() => setConfirmLogout(true)}
    className="flex items-center justify-center gap-2 mx-auto bg-red-600 text-white px-6 py-2 rounded-md font-semibold hover:bg-red-700 active:scale-95 transition"
  >
    <LogOut size={18} /> Logout
  </button>
  <p className="text-xs text-gray-400 mt-2">
    Logging out will clear your session and cookies.
  </p>

  {/* 🔴 Logout Confirmation Modal */}
  <AnimatePresence>
    {confirmLogout && (
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        exit={{ opacity: 0 }}
        className="fixed inset-0 bg-black/40 backdrop-blur-sm flex items-center justify-center z-50"
      >
        <motion.div
          initial={{ scale: 0.9, opacity: 0 }}
          animate={{ scale: 1, opacity: 1 }}
          exit={{ scale: 0.9, opacity: 0 }}
          className="bg-white dark:bg-neutral-900 border-2 border-black dark:border-white rounded-xl shadow-lg max-w-sm w-full p-6 text-center relative"
        >
          <button
            className="absolute top-3 right-3 text-gray-500 hover:text-black dark:hover:text-white"
            onClick={() => setConfirmLogout(false)}
          >
            <X size={20} />
          </button>

          <h2 className="text-xl font-bold mb-2 text-neutral-900 dark:text-neutral-100">
            Log Out?
          </h2>
          <p className="text-gray-600 dark:text-gray-400 mb-6">
            Are you sure you want to log out of your account?
          </p>

          <div className="flex justify-center gap-4">
            <button
              onClick={() => setConfirmLogout(false)}
              className="px-4 py-2 bg-gray-300 hover:bg-gray-400 rounded font-semibold transition"
            >
              Cancel
            </button>
            <button
              onClick={handleLogout}
              className="px-4 py-2 bg-red-600 text-white hover:bg-red-700 rounded font-semibold transition"
            >
              Logout
            </button>
          </div>
        </motion.div>
      </motion.div>
    )}
  </AnimatePresence>
</div>

      </motion.div>
    </div>
  );
}
