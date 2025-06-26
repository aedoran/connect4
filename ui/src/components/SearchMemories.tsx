import { gql, useLazyQuery } from '@apollo/client';
import { useState } from 'react';
import { MemoryDetails } from './MemoryDetails';

const SEARCH = gql`
  query Search($vector: [Float!], $limit: Int) {
    search(vector: $vector, limit: $limit) {
      id
      score
    }
  }
`;

interface Result {
  id: number;
  score: number;
}

export function SearchMemories() {
  const [vectorInput, setVectorInput] = useState('');
  const [limit, setLimit] = useState(5);
  const [selected, setSelected] = useState<Result | null>(null);
  const [runSearch, { data, loading }] = useLazyQuery<{ search: Result[] }>(SEARCH);

  const results = data?.search ?? [];

  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const vector = vectorInput
      .split(',')
      .map((v) => parseFloat(v.trim()))
      .filter((v) => !Number.isNaN(v));
    runSearch({ variables: { vector, limit } });
  };

  return (
    <div className="flex gap-4">
      <div className="flex-1">
        <form onSubmit={onSubmit} className="mb-4 flex gap-2">
          <input
            className="border px-2 py-1 flex-1"
            value={vectorInput}
            onChange={(e) => setVectorInput(e.target.value)}
            placeholder="Vector e.g. 0.1,0.2"
          />
          <input
            type="number"
            className="border px-2 py-1 w-16"
            value={limit}
            onChange={(e) => setLimit(parseInt(e.target.value, 10))}
          />
          <button className="px-4 py-1 bg-primary text-primary-foreground rounded">Search</button>
        </form>
        <ul>
          {loading && <li>Loading...</li>}
          {results.map((r) => (
            <li key={r.id}>
              <button onClick={() => setSelected(r)} className="underline text-blue-600">
                Memory {r.id} (score {r.score.toFixed(2)})
              </button>
            </li>
          ))}
        </ul>
      </div>
      {selected && (
        <div className="w-64 p-4 border-l">
          <h2 className="font-bold mb-2">Memory {selected.id}</h2>
          <MemoryDetails id={selected.id} />
        </div>
      )}
    </div>
  );
}
