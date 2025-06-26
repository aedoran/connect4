import { useCounterStore } from '../store';

export function Home() {
  const count = useCounterStore((s) => s.count);
  const increment = useCounterStore((s) => s.increment);
  return (
    <div>
      <h1 className="text-2xl font-bold mb-4">Home</h1>
      <button className="px-4 py-2 bg-primary text-primary-foreground rounded" onClick={increment}>
        Count is {count}
      </button>
    </div>
  );
}
