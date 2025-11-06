import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

const PUBLIC_ROUTES = ["/login", "/signup"];
const PROTECTED_PREFIX = "/(protected)";

// Cookie name from env or fallback
const COOKIE_NAME = process.env.COOKIE_NAME || "hl_jwt";

export function middleware(req: NextRequest) {
  const { pathname } = req.nextUrl;
  const token = req.cookies.get(COOKIE_NAME)?.value;

  const isPublic = PUBLIC_ROUTES.includes(pathname);
  const isProtected = pathname.startsWith("/(protected)");

  // 1. If user is logged in & visiting public route -> redirect to dashboard
  if (token && isPublic) {
    const url = req.nextUrl.clone();
    url.pathname = "/dashboard";
    return NextResponse.redirect(url);
  }

  // 2. If user is NOT logged in & trying to access protected -> redirect to login
  if (!token && isProtected) {
    const url = req.nextUrl.clone();
    url.pathname = "/login";
    return NextResponse.redirect(url);
  }

  // 3. Otherwise allow
  return NextResponse.next();
}

export const config = {
  matcher: [
    "/login",
    "/signup",
    "/(protected)/:path*",
  ],
};
