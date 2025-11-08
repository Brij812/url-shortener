"use client";

import { usePathname } from "next/navigation";
import Link from "next/link";

export default function ProtectedLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const pathname = usePathname();

  const navItems = [
    { href: "/dashboard", label: "Dashboard" },
    { href: "/create", label: "Create URL" },
    { href: "/metrics", label: "Metrics" },
    { href: "/settings", label: "Settings" },
  ];

  return (
    <div className="flex min-h-screen">
      {/* Sidebar */}
      <aside className="w-64 bg-white border-r border-gray-200 p-6 flex flex-col shadow-sm">
        <h1 className="text-2xl font-bold mb-8 tracking-tight">HyperLinkOS</h1>

        <nav className="space-y-2">
          {navItems.map((item) => (
            <Link
              key={item.href}
              href={item.href}
              className={`block px-3 py-2 rounded-md transition font-medium ${
                pathname === item.href
                  ? "bg-gray-900 text-white"
                  : "hover:bg-gray-100 text-gray-700"
              }`}
            >
              {item.label}
            </Link>
          ))}
        </nav>

        <div className="mt-auto pt-10 text-sm text-gray-400">
          © {new Date().getFullYear()} HyperLinkOS
        </div>
      </aside>

      {/* Main Content */}
      <main className="flex-1 p-10 bg-[#F7F7F8] animate-in fade-in slide-in-from-bottom-4 duration-300">
        {children}
      </main>
    </div>
  );
}
