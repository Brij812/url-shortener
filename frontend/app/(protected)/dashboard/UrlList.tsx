"use client";

import { Card } from "@/components/ui/card";

export function UrlList({ urls }: { urls: any[] }) {
  if (urls.length === 0) {
    return <p className="text-gray-500">You haven’t shortened any URLs yet.</p>;
  }

  return (
    <div className="space-y-4">
      {urls.map((u) => (
        <Card key={u.id} className="p-4 border">
          <div className="flex justify-between items-center">
            <div>
              <p className="font-medium">{u.long_url}</p>
              <p className="text-sm text-gray-500">{u.code}</p>
            </div>

            <a
              href={`${process.env.NEXT_PUBLIC_API_BASE_URL}/${u.code}`}
              target="_blank"
              className="text-blue-600 hover:underline"
            >
              Open
            </a>
          </div>
        </Card>
      ))}
    </div>
  );
}
