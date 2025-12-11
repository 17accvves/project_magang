import React, { useEffect, useState } from "react";
import { AiOutlineCheck, AiOutlineClose } from "react-icons/ai";
import "./AdminCafeApproval.css";

export default function AdminCafeApproval() {
  const [cafes, setCafes] = useState([]); // selalu array
  const [loading, setLoading] = useState(false);

  // ==============================
  // Fetch semua cafe dari backend
  // ==============================
  const fetchCafes = async () => {
    setLoading(true);
    try {
      const res = await fetch("http://localhost:8080/all-cafes");
      if (!res.ok) throw new Error("Gagal fetch cafes");

      const data = await res.json();

      // Pastikan data selalu array
      const arr = Array.isArray(data) ? data : [];
      const cafesWithStatus = arr.map((c) => ({
        ...c,
        status: c.verified ? "approved" : c.rejected ? "rejected" : "pending",
      }));

      setCafes(cafesWithStatus);
    } catch (err) {
      console.error(err);
      alert("Gagal mengambil data cafe");
      setCafes([]); // default agar tidak null
    } finally {
      setLoading(false);
    }
  };

  // ==============================
  // Approve / Tolak cafe
  // ==============================
  const updateCafeStatus = async (id, approve = true) => {
    try {
      const endpoint = approve
        ? "http://localhost:8080/approve-cafe"
        : "http://localhost:8080/reject-cafe";

      const res = await fetch(endpoint, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ cafe_id: id }),
      });

      if (!res.ok) throw new Error("Gagal update cafe");

      const data = await res.json();
      alert(data.message);

      // Update status di state tanpa menghapus baris
      setCafes((prev) =>
        prev.map((c) =>
          c.id === id ? { ...c, status: approve ? "approved" : "rejected" } : c
        )
      );
    } catch (err) {
      console.error(err);
      alert("Gagal update status cafe");
    }
  };

  useEffect(() => {
    fetchCafes();
  }, []);

  return (
    <div className="admin-container" style={{ padding: "20px" }}>
      <h2>Daftar Cafe</h2>

      {loading ? (
        <p>Loading...</p>
      ) : cafes.length === 0 ? (
        <p>Tidak ada cafe</p>
      ) : (
        <table
          border="1"
          cellPadding="10"
          style={{ borderCollapse: "collapse", width: "100%" }}
        >
          <thead>
            <tr>
              <th>Username</th>
              <th>Email</th>
              <th>Izin Usaha</th>
              <th>Status</th>
              <th>Action</th>
            </tr>
          </thead>
          <tbody>
            {cafes.map((cafe) => (
              <tr key={cafe.id}>
                <td>{cafe.username}</td>
                <td>{cafe.email}</td>
                <td>
                  <a
                    href={`http://localhost:8080/${cafe.izin_usaha}`}
                    target="_blank"
                    rel="noreferrer"
                  >
                    Lihat File
                  </a>
                </td>
                <td>
                  {cafe.status === "pending"
                    ? "Menunggu"
                    : cafe.status === "approved"
                    ? "Disetujui"
                    : "Ditolak"}
                </td>
                <td>
                  {cafe.status === "pending" ? (
                    <>
                      <button
                        onClick={() => updateCafeStatus(cafe.id, true)}
                        style={{
                          backgroundColor: "#4CAF50",
                          color: "white",
                          padding: "5px 10px",
                          border: "none",
                          borderRadius: "4px",
                          cursor: "pointer",
                          marginRight: "5px",
                        }}
                      >
                        <AiOutlineCheck /> Approve
                      </button>
                      <button
                        onClick={() => updateCafeStatus(cafe.id, false)}
                        style={{
                          backgroundColor: "#f44336",
                          color: "white",
                          padding: "5px 10px",
                          border: "none",
                          borderRadius: "4px",
                          cursor: "pointer",
                        }}
                      >
                        <AiOutlineClose /> Tolak
                      </button>
                    </>
                  ) : (
                    <span
                      style={{
                        color:
                          cafe.status === "approved" ? "#4CAF50" : "#f44336",
                        fontWeight: "bold",
                      }}
                    >
                      {cafe.status === "approved" ? "Disetujui" : "Ditolak"}
                    </span>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}
