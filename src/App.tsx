import React, { useState, useEffect, useCallback, KeyboardEvent } from 'react';
import axios from 'axios';

interface SearchResult {
  results: string[] | null;
  exists: boolean;
}

function App(): JSX.Element {
  const [query, setQuery] = useState<string>('');
  const [results, setResults] = useState<string[] | null>(null);

  const handleSearch = useCallback(async () => {
    try {
      const response = await axios.post<SearchResult>('http://localhost:3000/search', { query });
      console.log('Response from server:', response.data);

      if (!response.data.exists) {
        console.log('Data not found in Redis');
        setResults(null);
        return;
      }

      setResults(response.data.results);
    } catch (error) {
      console.error('Error:', error);
    }
  }, [query]);

  useEffect(() => {
    handleSearch();
  }, [handleSearch]);

  const handleKeyPress = (event: KeyboardEvent<HTMLInputElement>): void => {
    // Trigger search when Enter key is pressed
    if (event.key === 'Enter') {
      handleSearch();
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-300 via-blue-200 to-blue-100">
      <div className="bg-white p-8 rounded shadow-md w-96 text-center animate__animated animate__fadeIn">
        <h1 className="text-3xl font-semibold mb-4 text-gray-800">Search Engine</h1>
        <div className="mx-auto mb-4 w-3/4">
          <input
            type="text"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            onKeyPress={handleKeyPress}
            className="border p-2 w-full focus:outline-none focus:ring focus:border-blue-500 transition duration-300"
            placeholder="Search..."
          />
        </div>
        {results ? (
          <div className="mt-4">
            {results.map((result, index) => (
              <div key={index} className="mb-2 text-gray-800">
                {result}
              </div>
            ))}
          </div>
        ) : (
          <p className="mt-4 text-gray-600">{query} - No results available</p>
        )}
      </div>
    </div>
  );
}

export default App;
