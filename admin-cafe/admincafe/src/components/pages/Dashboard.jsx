import React from "react";
import { Line } from "react-chartjs-2";
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
} from "chart.js";
import { FaUser, FaUtensils, FaComment, FaShoppingBasket } from "react-icons/fa";
import "./Dashboard.css";

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend
);

function Dashboard() {
  const salesData = {
    labels: Array.from({ length: 30 }, (_, i) => i),
    datasets: [
      {
        label: "Penjualan",
        data: [
          1900000, 800000, 500000, 350000, 1100000, 1200000, 500000, 500000,
          1700000, 1800000, 500000, 100000, 200000, 300000, 650000, 1100000,
          2000000, 1200000, 450000, 1450000, 1700000, 350000, 1050000, 750000,
          300000, 350000, 2000000, 600000, 2000000, 1800000,
        ],
        borderColor: "#4F46E5",
        backgroundColor: "rgba(79,70,229,0.2)",
        tension: 0.3,
      },
    ],
  };

  const salesOptions = {
    responsive: true,
    plugins: {
      legend: { display: false },
      title: { display: true, text: "Grafik Penjualan" },
    },
    scales: {
      y: {
        ticks: {
          callback: (value) => `Rp. ${value.toLocaleString()}`,
        },
      },
    },
  };

  return (
    <div className="dashboard">
      {/* Stats Cards */}
      <div className="stats-cards">
        <div className="card">
          <div className="card-icon"><FaUser /></div>
          <div className="card-info">
            <p>Jumlah Pesanan</p>
            <h3>120</h3>
          </div>
        </div>
        <div className="card">
          <div className="card-icon"><FaUtensils /></div>
          <div className="card-info">
            <p>Jumlah Menu</p>
            <h3>005</h3>
          </div>
        </div>
        <div className="card">
          <div className="card-icon"><FaComment /></div>
          <div className="card-info">
            <p>Rating</p>
            <h3>BINTANG ?/5</h3>
          </div>
        </div>
        <div className="card">
          <div className="card-icon"><FaShoppingBasket /></div>
          <div className="card-info">
            <p>Menu Terlaris</p>
            <ol>
              <li>Extra Jose susu</li>
              <li>Kerupuk</li>
              <li>Kopi</li>
            </ol>
          </div>
        </div>
      </div>

      {/* Chart */}
      <div className="chart-card">
        <Line data={salesData} options={salesOptions} />
      </div>

      {/* Bottom Cards */}
      <div className="bottom-cards">
        <div className="card event-card">
          <h4>Event Mendatang</h4>
          <p className="emoji">📌 Acoustic Night</p>
          <p className="emoji">📅 30 Sep 2025, 19.00 - 22.00</p>
          <p className="emoji">📍 Cafe ABC, Jakarta</p>
          <p className="emoji">👥 Kapasitas: 50 orang</p>
          <p className="emoji">🎫 Tiket: Gratis</p>
        </div>
        <div className="card promo-card">
          <h4>Promo Aktif</h4>
          <ul>
            <li>Diskon Matcha 25%</li>
            <li>Diskon Pentolan Pey 35%</li>
          </ul>
        </div>
        <div className="card review-card">
          <h4>Review Terbaru</h4>
          <p><b>James:</b> "Kopinya enak banget!" ⭐⭐⭐⭐⭐</p>
          <p><b>Rodri:</b> "Suasana tempatnya adem!" ⭐⭐⭐⭐</p>
        </div>
      </div>
    </div>
  );
}

export default Dashboard;
