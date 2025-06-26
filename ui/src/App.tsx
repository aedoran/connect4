import { Route, Routes, Link } from 'react-router-dom';
import { Home } from './pages/Home';
import { About } from './pages/About';

export default function App() {
  return (
    <div className="container mx-auto p-4">
      <nav className="mb-4 flex gap-4">
        <Link to="/" className="font-bold">Home</Link>
        <Link to="/about">About</Link>
      </nav>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/about" element={<About />} />
      </Routes>
    </div>
  );
}
