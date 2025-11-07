import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

const PUBLIC_ROUTES = ["/login", "/signup"];
const COOKIE_NAME = process.env.COOKIE_NAME || "hl_jwt";

export function middleware(req: NextRequest) {
  const { pathname } = req.nextUrl;
  const token = req.cookies.get(COOKIE_NAME)?.value;

  const isPublic = PUBLIC_ROUTES.includes(pathname);
  const isProtected =
    pathname.startsWith("/dashboard") ||
    pathname.startsWith("/create") ||
    pathname.startsWith("/metrics") ||
    pathname.startsWith("/settings");

  // logged in user visiting login/signup → redirect dashboard
  if (token && isPublic) {
    const url = req.nextUrl.clone();
    url.pathname = "/dashboard";
    return NextResponse.redirect(url);
  }

  // not logged in user visiting protected → redirect login
  if (!token && isProtected) {
    const url = req.nextUrl.clone();
    url.pathname = "/login";
    return NextResponse.redirect(url);
  }

  return NextResponse.next();
}

export const config = {
  matcher: [
    "/login",
    "/signup",
    "/dashboard/:path*",
    "/create/:path*",
    "/metrics/:path*",
    "/settings/:path*",
  ],
};
