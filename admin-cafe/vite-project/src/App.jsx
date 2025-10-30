import { BrowserRouter as Router, Routes, Route, Navigate } from "react-router-dom";
import Login from "./Login";
import Register from "./Register";
import Admincafe from "./components/Admincafe";
import React from "react";

function App() {
  return (
    <Router>
      <Routes>
        {/* 🔹 Route default (redirect ke /login) */}
        <Route path="/" element={<Navigate to="/login" replace />} />
        
        {/* 🔹 Halaman login dan register */}
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/admin" element={<Admincafe />} />
      </Routes>
    </Router>
  );
}

export default App;
