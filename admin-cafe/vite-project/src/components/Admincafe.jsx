import React, { useState } from "react";
import {
  Home,
  User,
  Coffee,
  Star,
  Tag,
  Calendar,
  BarChart2,
  LogOut,
  Bell,
  ChevronLeft,
  ChevronRight,
} from "lucide-react";
import { useNavigate } from "react-router-dom";
import "./Admincafe.css";
import coffeeLogo from "../assets/gambarbiji.png";
import profilePic from "../assets/profile.jpg";

// === IMPORT HALAMAN ===
import Dashboard from "./pages/Dashboard";
import ProfileCafe from "./pages/ProfileCafe";
import MenuPage from "./pages/Menu";
import UlasanPage from "./pages/Ulasan";
import PromoPage from "./pages/Promo";
import EventPage from "./pages/Event";
import LaporanPage from "./pages/Laporan";

function Admincafe() {
  const [activeMenu, setActiveMenu] = useState("Dashboard");
  const [isCollapsed, setIsCollapsed] = useState(false);
  const navigate = useNavigate();

  const menuItems = [
    { name: "Dashboard", icon: <Home size={20} /> },
    { name: "Profile Cafe", icon: <User size={20} /> },
    { name: "Menu", icon: <Coffee size={20} /> },
    { name: "Ulasan", icon: <Star size={20} /> },
    { name: "Promo", icon: <Tag size={20} /> },
    { name: "Event", icon: <Calendar size={20} /> },
    { name: "Laporan", icon: <BarChart2 size={20} /> },
  ];

  const handleLogout = () => {
    const confirmLogout = window.confirm("Yakin ingin keluar dari Admin Café?");
    if (confirmLogout) {
      localStorage.removeItem("isLoggedIn");
      navigate("/login");
    }
  };

  // === FUNGSI RENDER HALAMAN ===
  const renderPage = () => {
    switch (activeMenu) {
      case "Dashboard":
        return <Dashboard />;
      case "Profile Cafe":
        return <ProfileCafe />;
      case "Menu":
        return <MenuPage />;
      case "Ulasan":
        return <UlasanPage />;
      case "Promo":
        return <PromoPage />;
      case "Event":
        return <EventPage />;
      case "Laporan":
        return <LaporanPage />;
      default:
        return <Dashboard />;
    }
  };

  return (
    <div className={`admin-layout ${isCollapsed ? "collapsed" : ""}`}>
      {/* === SIDEBAR === */}
      <aside className="sidebar">
        <div className="logo-section">
          <img src={coffeeLogo} alt="Coffee" className="logo-icon" />
          {!isCollapsed && <h1 className="logo-text">Carispot</h1>}
        </div>

        <button
          className="toggle-btn"
          onClick={() => setIsCollapsed(!isCollapsed)}
        >
          {isCollapsed ? <ChevronRight size={20} /> : <ChevronLeft size={20} />}
        </button>

        <nav className="menu">
          {menuItems.map((item) => (
            <div
              key={item.name}
              className={`menu-item ${
                activeMenu === item.name ? "active" : ""
              }`}
              onClick={() => setActiveMenu(item.name)}
            >
              {item.icon}
              {!isCollapsed && <span>{item.name}</span>}
            </div>
          ))}
        </nav>

        <div className="logout" onClick={handleLogout}>
          <LogOut size={20} />
          {!isCollapsed && <span>Logout</span>}
        </div>
      </aside>

      {/* === KONTEN UTAMA === */}
      <main className="main-content">
        <header className="top-bar-admin">
          <div className="top-bar-left">{activeMenu}</div>
          <div className="top-bar-right">
            <Bell className="notif-icon" size={22} />
            <div className="profile-info">
              <img src={profilePic} alt="Profile" className="profile-pic" />
              <span className="profile-name">Admin Café</span>
            </div>
          </div>
        </header>

        <div className="content-body">
          {renderPage()}
        </div>
      </main>
    </div>
  );
}

export default Admincafe;
