import React from "react";
import { MapContainer, TileLayer, Marker, useMap } from "react-leaflet";
import "leaflet/dist/leaflet.css";

function getPosition() {
  return new Promise((resolve, reject) => {
    navigator.geolocation.getCurrentPosition(resolve, reject);
  });
}
export default function Home() {
  navigator.geolocation.getCurrentPosition((position) => {
    alert(position.coords.latitude, position.coords.longitude);
  });
  return (
    <div>
      <MapContainer
        center={[15, 15]}
        zoom={15}
        style={{ height: "500px", width: "100%" }}
      >
        <TileLayer url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png" />
        {/* <Marker position={[position.lat, position.lng]} />
        <RecenterMap position={position} /> */}
      </MapContainer>
    </div>
  );
}
