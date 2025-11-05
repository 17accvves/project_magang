import React, { useState, useMemo } from "react";
import "./Promo.css";
import { FaChevronUp, FaChevronDown } from "react-icons/fa";

const initialPromos = [
  {
    id: 1,
    title: "Diskon Promo 25% Matcha",
    start: "18 September 2025",
    end: "25 September 2025",
    revenue: "17JT",
  },
  {
    id: 2,
    title: "Diskon Promo 15% Kopi Latte",
    start: "10 September 2025",
    end: "20 September 2025",
    revenue: "12JT",
  },
  {
    id: 3,
    title: "Diskon Promo 10% Snack Pukis",
    start: "5 September 2025",
    end: "12 September 2025",
    revenue: "9JT",
  },
  {
    id: 4,
    title: "Diskon Promo 20% Extra Jose",
    start: "1 September 2025",
    end: "10 September 2025",
    revenue: "14JT",
  },
  {
    id: 5,
    title: "Diskon Promo 30% Pentolaan James",
    start: "22 Agustus 2025",
    end: "28 Agustus 2025",
    revenue: "20JT",
  },
];

const Promo = () => {
  const [promos] = useState(initialPromos);
  const [showScrollTop, setShowScrollTop] = useState(true);

  // Hitung total revenue (anggap “JT” = juta)
  const totalRevenue = useMemo(() => {
    return promos.reduce((sum, promo) => {
      const number = parseInt(promo.revenue.replace("JT", "").trim());
      return sum + (isNaN(number) ? 0 : number);
    }, 0);
  }, [promos]);

  const toggleScroll = () => {
    setShowScrollTop(!showScrollTop);
    window.scrollTo({
      top: showScrollTop ? 0 : document.body.scrollHeight,
      behavior: "smooth",
    });
  };

  return (
    <div className="promo-container">
      {/* Sidebar kiri */}
      <div className="promo-sidebar">
        <p><strong>Promo Aktif :</strong> {promos.length}</p>
        <p><strong>Total Revenue :</strong> {totalRevenue}JT</p>
        <p><strong>Daftar Promo Aktif :</strong></p>
        <ul>
          {promos.map((promo) => (
            <li key={promo.id}>{promo.title}</li>
          ))}
        </ul>
      </div>

      {/* Grid promo */}
      <div className="promo-content">
        <div className="promo-grid">
          {promos.map((promo) => (
            <div key={promo.id} className="promo-card">
              <h4 className="promo-title">{promo.title}</h4>

              <table className="promo-table">
                <thead>
                  <tr>
                    <th>Mulai Promo</th>
                    <th>Akhir Promo</th>
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td>{promo.start}</td>
                    <td>{promo.end}</td>
                  </tr>
                </tbody>
              </table>

              <p className="promo-revenue">REVENUE : {promo.revenue}</p>
            </div>
          ))}
        </div>

        {/* Scroll icon */}
        <div className="scroll-icon" onClick={toggleScroll}>
          {showScrollTop ? <FaChevronUp /> : <FaChevronDown />}
        </div>
      </div>
    </div>
  );
};

export default Promo;
