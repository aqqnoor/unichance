import { Routes, Route } from 'react-router-dom';
import { useState } from 'react';
import { StudentProfile } from './types';

import Home from './pages/Home';
import Results from './pages/Results';
import Profile from './pages/Profile';
import Search from './pages/Search';
import University from "./pages/University";


import Layout from './components/Layout';

function App() {
  const [profile, setProfile] = useState<StudentProfile>({});

  return (
    <Routes>
      {/* Layout — parent */}
      <Route element={<Layout />}>
        <Route
          path="/"
          element={<Home profile={profile} setProfile={setProfile} />}
        />
        <Route
          path="/results"
          element={<Results profile={profile} />}
        />
        <Route
          path="/profile"
          element={<Profile profile={profile} setProfile={setProfile} />}
        />
        <Route path="/search" element={<Search />} />
        <Route path="/universities/:id" element={<University />} />

      </Route>
    </Routes>
  );
}

export default App;
