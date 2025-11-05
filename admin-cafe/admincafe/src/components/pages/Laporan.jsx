import React, { useState } from "react";
import "./Laporan.css";
import { Line } from "react-chartjs-2";
import { FaStar, FaHome, FaUtensils } from "react-icons/fa";
import "chart.js/auto";

const Laporan = () => {
  const [startDate, setStartDate] = useState("2025-09-01");
  const [endDate, setEndDate] = useState("2025-09-10");

  // Data simulasi (nanti bisa diganti dari backend)
  const dataPenjualan = [100000, 600000, 900000, 850000, 920000, 600000, 450000, 550000, 300000, 100000];
  const tanggal = ["1", "2", "3", "4", "5", "6", "7", "8", "9", "10"];

  const totalPenjualan = "Rp. 12.023.000";
  const totalPesanan = 370;
  const ratingRata = "4.6/5";
  const menuTerlaris = ["Extra Jose Susu", "Pukis Pey", "Pentolan James"];

  const chartData = {
    labels: tanggal,
    datasets: [
      {
        label: "Pendapatan Harian",
        data: dataPenjualan,
        fill: false,
        borderColor: "#b57b3f",
        backgroundColor: "#b57b3f",
        tension: 0.4,
      },
    ],
  };

  const chartOptions = {
    plugins: { legend: { display: false } },
    scales: {
      y: {
        ticks: { callback: (value) => `Rp. ${value.toLocaleString()}` },
      },
    },
  };

  const daftarMenu = [
    { tanggal: "1 September 2025", menu: "Extra Jose Susu, Pukis pey, Pentolan James", terjual: 67, pendapatan: "Rp. 9.000.000" },
    { tanggal: "2 September 2025", menu: "Extra Jose Susu, Pukis pey, Pentolan James", terjual: 67, pendapatan: "Rp. 7.400.000" },
    { tanggal: "3 September 2025", menu: "Extra Jose Susu, Pukis pey, Pentolan James", terjual: 67, pendapatan: "Rp. 5.890.000" },
  ];

  return (
    <div className="laporan-container">
      <h1 className="laporan-title">Laporan</h1>

      <div className="laporan-filter">
        <label>
          Tanggal mulai :
          <input
            type="date"
            value={startDate}
            onChange={(e) => setStartDate(e.target.value)}
          />
        </label>
        <span>To</span>
        <label>
          Tanggal akhir :
          <input
            type="date"
            value={endDate}
            onChange={(e) => setEndDate(e.target.value)}
          />
        </label>
        <button className="btn-oke">Oke</button>

        <div className="export-section">
          <span>Export :</span>
          <button className="export-btn">Pdf</button>
          <button className="export-btn">Excel</button>
        </div>
      </div>

      {/* Kartu statistik */}
      <div className="laporan-cards">
        <div className="laporan-card">
          <div className="icon-container"><FaHome /></div>
          <div>
            <p className="card-title">Total Penjualan</p>
            <h3>{totalPenjualan}</h3>
          </div>
        </div>

        <div className="laporan-card">
          <div className="icon-container"><FaUtensils /></div>
          <div>
            <p className="card-title">Total Pemesanan</p>
            <h3>{totalPesanan}</h3>
          </div>
        </div>

        <div className="laporan-card">
          <div className="icon-container star"><FaStar /></div>
          <div>
            <p className="card-title">Rating Rata-rata</p>
            <h3>{ratingRata}</h3>
          </div>
        </div>

        <div className="laporan-card">
          <div>
            <p className="card-title">Menu Terlaris</p>
            <ol>
              {menuTerlaris.map((item, i) => (
                <li key={i}>{i + 1}. {item}</li>
              ))}
            </ol>
          </div>
        </div>
      </div>

      {/* Grafik */}
      <div className="laporan-chart">
        <Line data={chartData} options={chartOptions} />
      </div>

      {/* Daftar Menu */}
      <div className="laporan-table-container">
        <h3 className="table-title">DAFTAR MENU</h3>
        <table className="laporan-table">
          <thead>
            <tr>
              <th>Tanggal</th>
              <th>Nama Menu</th>
              <th>Terjual</th>
              <th>Pendapatan</th>
            </tr>
          </thead>
          <tbody>
            {daftarMenu.map((item, index) => (
              <tr key={index}>
                <td>{item.tanggal}</td>
                <td>{item.menu}</td>
                <td>{item.terjual}</td>
                <td>{item.pendapatan}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default Laporan;
