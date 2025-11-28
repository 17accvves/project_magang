import React, { useState } from "react";
import "./Event.css";
import { FaCheck, FaTimes } from "react-icons/fa";
import CircularText from './item/CircularText';

const eventData = {
  request: [
    { pengaju: "guest113", judul: "Workshop Latte Art", mulai: "10 Sep 2025", selesai: "19 Sep 2025" },
    { pengaju: "guest223", judul: "Workshop Latte Art", mulai: "10 Sep 2025", selesai: "19 Sep 2025" },
    { pengaju: "SMAN 1 KILO", judul: "Workshop Latte Art", mulai: "10 Sep 2025", selesai: "19 Sep 2025" },
  ],
  berlangsung: [
    { pengaju: "eventanakmuda", judul: "Cafe Night Gathering", mulai: "01 Nov 2025", selesai: "05 Nov 2025" },
    { pengaju: "guest312", judul: "Mini Coffee Talk", mulai: "02 Nov 2025", selesai: "06 Nov 2025" },
  ],
  ditolak: [
    { pengaju: "guest333", judul: "Event Kopi Gratis", mulai: "10 Sep 2025", selesai: "19 Sep 2025" },
    { pengaju: "guest114", judul: "Ngopi dan Nulis", mulai: "5 Okt 2025", selesai: "8 Okt 2025" },
  ],
  selesai: [
    { pengaju: "SMAN 1 KILO", judul: "Workshop Latte Art", mulai: "10 Sep 2025", selesai: "19 Sep 2025" },
    { pengaju: "guest223", judul: "Kelas Latte Lanjutan", mulai: "15 Sep 2025", selesai: "20 Sep 2025" },
  ],
};

const Event = () => {
  const [selectedTab, setSelectedTab] = useState("request");

  const handleTabClick = (tab) => {
    setSelectedTab(tab);
  };

  const data = eventData[selectedTab];

  return (
    <div className="event-container">
      <h2 className="event-title">Manajemen Event ({selectedTab.charAt(0).toUpperCase() + selectedTab.slice(1)})</h2>

      <div className="event-tab-container">
        {["request", "berlangsung", "ditolak", "selesai"].map((tab) => (
          <button
            key={tab}
            onClick={() => handleTabClick(tab)}
            className={`event-tab ${selectedTab === tab ? "active" : ""}`}
          >
            {tab === "request" && "Permintaan"}
            {tab === "berlangsung" && "Berlangsung"}
            {tab === "ditolak" && "Ditolak"}
            {tab === "selesai" && "Selesai"}
          </button>
        ))}
      </div>

      <div className="event-table-container">
        <table className="event-table">
          <thead>
            <tr>
              <th>Pengaju</th>
              <th>Judul Event</th>
              <th>Tanggall Mulai</th>
              <th>Tanggal Selesai</th>
              <th>Daftar Menu</th>
              <th>Persetujuan</th>
            </tr>
          </thead>
          <tbody>
            {data.map((item, index) => (
              <tr key={index}>
                <td>{item.pengaju}</td>
                <td>{item.judul}</td>
                <td>{item.mulai}</td>
                <td>{item.selesai}</td>
                <td>
                  <select className="menu-select">
                    <option>Menu</option>
                    <option>Kopi Latte</option>
                    <option>Matcha</option>
                    <option>Snack</option>
                  </select>
                </td>
                <td>
                  <div className="approval-buttons">
                    <button className="approve-btn">
                      <FaCheck />
                    </button>
                    <button className="reject-btn">
                      <FaTimes />
                    </button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <div className="event-pagination">
        <span>1-20 of 300 Row/Page:</span>
        <select>
          <option>7/12</option>
          <option>10/20</option>
        </select>
        <div className="pagination-btns">
          <button>{"<<"}</button>
          <button>{"<"}</button>
          <button className="active">1</button>
          <button>2</button>
          <button>3</button>
          <button>{">"}</button>
          <button>{">>"}</button>
        </div>
      </div>

  
<CircularText
  text="Menu"
  onHover="speedUp"
  spinDuration={20}
  className="custom-class"
/>
    </div>
  );
};

export default Event;
