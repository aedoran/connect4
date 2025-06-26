import { useEffect, useState } from 'react';

interface Memory {
  id: number;
  userID: number;
  content: string;
  createdAt: string;
}

export function MemoryDetails({ id }: { id: number }) {
  const [memory, setMemory] = useState<Memory | null>(null);

  useEffect(() => {
    fetch(`${import.meta.env.VITE_API_URL}/api/v1/memories/${id}`)
      .then((r) => r.json())
      .then((d) => setMemory(d));
  }, [id]);

  if (!memory) {
    return <div>Loading...</div>;
  }

  return (
    <div>
      <p className="mb-2">{memory.content}</p>
      <p className="text-sm text-gray-500">User {memory.userID}</p>
    </div>
  );
}
