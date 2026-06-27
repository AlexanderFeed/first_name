import React from 'react';
import { BrowserRouter, Routes, Route, Link, Outlet } from 'react-router-dom';
import ReactDOM from 'react-dom/client';
import Register from './register';
import Login from './login';

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <BrowserRouter>
    <nav>
      <Link to="/">Home</Link>
      <Link to="/login">Login</Link>
      <Link to="/register">Register</Link>
    </nav>
    <Routes>
      {/* <Route path="/" element={< />}> */}
      <Route path="/register" element={< Register/>}/>
      <Route path="/login" element={< Login/>}/>
    </Routes>
    </BrowserRouter>
  </React.StrictMode>
);

