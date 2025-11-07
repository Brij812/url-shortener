import { cookies } from "next/headers";
import axios from "axios";
import { UrlList } from "./UrlList";

export const metadata = {
  title: "Dashboard | HyperLinkOS",
};

export default async function DashboardPage() {
  const cookieStore = await cookies();
  const token = cookieStore.get(process.env.COOKIE_NAME || "hl_jwt")?.value;

  if (!token) {
    return <div className="p-8">Not authenticated</div>;
  }

  const res = await axios.get(
    `${process.env.NEXT_PUBLIC_API_BASE_URL}/all`,
    {
      headers: {
        Authorization: `Bearer ${token}`,
      },
      validateStatus: () => true,
    }
  );

  const urls = res.data || [];

  return (
    <div className="max-w-4xl mx-auto mt-10">
      <h1 className="text-3xl font-semibold mb-6">Your Shortened URLs</h1>
      <UrlList urls={urls} />
    </div>
  );
}
