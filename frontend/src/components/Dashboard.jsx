import React, { useState, useEffect, useRef } from 'react';
import { motion, AnimatePresence } from 'framer-motion';

const Dashboard = () => {
  const [data, setData] = useState({
    warehouse: { Name: "Central Medical Depot", Inventory: { "Antibiotics": 5000, "Painkillers": 10000, "Vaccines": 2000, "Bandages": 5000 } },
    outlets: [
      { id: "Outlet-1", Name: "City General Hospital", Inventory: { "Antibiotics": 50, "Painkillers": 100, "Vaccines": 20, "Bandages": 200 } },
      { id: "Outlet-2", Name: "Community Clinic South", Inventory: { "Antibiotics": 50, "Painkillers": 100, "Vaccines": 20, "Bandages": 200 } },
      { id: "Outlet-3", Name: "University Medical Center", Inventory: { "Antibiotics": 50, "Painkillers": 100, "Vaccines": 20, "Bandages": 200 } },
      { id: "Outlet-4", Name: "Metro Pharmacy", Inventory: { "Antibiotics": 50, "Painkillers": 100, "Vaccines": 20, "Bandages": 200 } }
    ],
    date: '2024-01-01'
  });
  const [outletLogs, setOutletLogs] = useState({
    'Outlet-1': [],
    'Outlet-2': [],
    'Outlet-3': [],
    'Outlet-4': []
  });
  const [shipments, setShipments] = useState([]);
  const [connected, setConnected] = useState(false);
  const ws = useRef(null);

  useEffect(() => {
    ws.current = new WebSocket('ws://localhost:8080/ws');

    ws.current.onopen = () => {
      console.log('WebSocket connected');
      setConnected(true);
    };

    ws.current.onmessage = (event) => {
      const message = JSON.parse(event.data);
      if (message.type === 'state') {
        setData(prev => ({
          ...prev,
          date: message.date,
          // Only update warehouse/outlets if provided in message
          warehouse: message.warehouse || prev.warehouse,
          outlets: message.outlets || prev.outlets,
        }));
      } else if (message.type === 'log') {
        // Parse log to determine which outlet it belongs to
        const logText = message.message;
        let outletId = null;

        // Check which outlet this log is for
        if (logText.includes('City General') || logText.includes('Outlet-1')) outletId = 'Outlet-1';
        else if (logText.includes('Community Clinic') || logText.includes('Outlet-2')) outletId = 'Outlet-2';
        else if (logText.includes('University Medical') || logText.includes('Outlet-3')) outletId = 'Outlet-3';
        else if (logText.includes('Metro Pharmacy') || logText.includes('Outlet-4')) outletId = 'Outlet-4';

        if (outletId) {
          setOutletLogs((prev) => ({
            ...prev,
            [outletId]: [logText, ...prev[outletId].slice(0, 9)]
          }));
        }
      } else if (message.type === 'shipment') {
        // Trigger shipment animation
        const newShipment = { ...message, id: Date.now() + Math.random() };
        setShipments((prev) => [...prev, newShipment]);
        // Remove after animation
        setTimeout(() => {
          setShipments((prev) => prev.filter(s => s.id !== newShipment.id));
        }, 20000);
      }
    };

    return () => ws.current.close();
  }, []);

  if (!connected) return <div className="text-white text-center mt-20">Connecting to Supply Chain OS... (Ensure Go Backend is running)</div>;

  return (
    <div className="min-h-screen bg-slate-900 text-white p-8 font-sans overflow-hidden">
      <header className="mb-8 flex justify-between items-center z-10 relative">
        <h1 className="text-3xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-blue-400 to-emerald-400">
          Responsive AI Clusters (Pharma Supply Chain)
        </h1>
        <div className="flex flex-col items-end">
          <div className="text-xl font-mono text-emerald-300">
            Date: {data.date}
          </div>
          <div className="text-xs bg-blue-900/50 px-2 py-1 rounded border border-blue-500/50 mt-1">
            PROTOCOL: <span className="font-bold text-white">A2A (JSON-RPC)</span> | AUTONOMOUS
          </div>
        </div>
      </header>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 relative h-[800px]">

        {/* Render Shipments */}
        <AnimatePresence>
          {shipments.map((shipment) => (
            <ShipmentPackage key={shipment.id} shipment={shipment} />
          ))}
        </AnimatePresence>

        {/* Top Left - Outlet 1 with Chat */}
        <div className="absolute top-0 left-0 w-80 z-10">
          <OutletCard outlet={data.outlets[0]} id="Outlet-1" />
          <ConversationBox outletId="Outlet-1" logs={outletLogs['Outlet-1']} name="City General Hospital" />
        </div>

        {/* Top Right - Outlet 2 with Chat */}
        <div className="absolute top-0 right-0 w-80 z-10">
          <OutletCard outlet={data.outlets[1]} id="Outlet-2" />
          <ConversationBox outletId="Outlet-2" logs={outletLogs['Outlet-2']} name="Community Clinic South" />
        </div>

        {/* Center Warehouse */}
        <div className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 z-10">
          <WarehouseCard warehouse={data.warehouse} />
        </div>

        {/* Bottom Left - Outlet 3 with Chat */}
        <div className="absolute bottom-0 left-0 w-80 z-10">
          <OutletCard outlet={data.outlets[2]} id="Outlet-3" />
          <ConversationBox outletId="Outlet-3" logs={outletLogs['Outlet-3']} name="University Medical Center" />
        </div>

        {/* Bottom Right - Outlet 4 with Chat */}
        <div className="absolute bottom-0 right-0 w-80 z-10">
          <OutletCard outlet={data.outlets[3]} id="Outlet-4" />
          <ConversationBox outletId="Outlet-4" logs={outletLogs['Outlet-4']} name="Metro Pharmacy" />
        </div>

        {/* Connecting Lines (Visual only) */}
        <svg className="absolute inset-0 w-full h-full pointer-events-none opacity-20">
          <line x1="20%" y1="10%" x2="50%" y2="50%" stroke="cyan" strokeWidth="2" strokeDasharray="5,5" />
          <line x1="80%" y1="10%" x2="50%" y2="50%" stroke="cyan" strokeWidth="2" strokeDasharray="5,5" />
          <line x1="20%" y1="90%" x2="50%" y2="50%" stroke="cyan" strokeWidth="2" strokeDasharray="5,5" />
          <line x1="80%" y1="90%" x2="50%" y2="50%" stroke="cyan" strokeWidth="2" strokeDasharray="5,5" />
        </svg>

      </div>
    </div>
  );
};

// Conversation Box Component for each outlet
const ConversationBox = ({ outletId, logs, name }) => {
  if (logs.length === 0) return null;

  const renderLog = (log) => {
    // Highlight specific phrases to visualize Protocol Logic
    if (log.includes("Analysis:")) return <span className="text-cyan-300 block border-l-2 border-cyan-500 pl-2">{log}</span>;
    if (log.includes("Reasoning:")) return <span className="text-yellow-200 block border-l-2 border-yellow-500 pl-2">{log}</span>;
    if (log.includes("Proposal:")) return <span className="text-yellow-400 block ml-2 font-bold">{log}</span>;
    if (log.includes("Agreement reached")) return <span className="text-green-400 block font-bold border-l-2 border-green-500 pl-2">{log}</span>;
    if (log.includes("Recall")) return <span className="text-red-400 block font-bold border-l-2 border-red-500 pl-2">{log}</span>;
    if (log.includes("Event:")) return <span className="text-purple-300 block font-bold mb-1">{log}</span>;
    if (log.includes("Request:")) return <span className="text-orange-300 block ml-2">{log}</span>;
    return <span className="text-emerald-300/80">{log}</span>;
  };

  return (
    <div className="mt-2 bg-black/90 p-3 rounded-lg border border-blue-500/30 max-h-64 overflow-y-auto shadow-2xl backdrop-blur-sm">
      <h4 className="text-xs font-bold text-blue-400 mb-2 flex items-center gap-2">
        <div className="w-2 h-2 rounded-full bg-green-500 animate-pulse"></div>
        A2A LOGS: {name}
      </h4>
      <div className="flex flex-col gap-3">
        {logs.map((log, i) => (
          <div key={i} className="text-[10px] md:text-xs font-mono whitespace-pre-wrap leading-relaxed">
            {renderLog(log)}
          </div>
        ))}
      </div>
    </div>
  );
};

// Moving Package Animation
const ShipmentPackage = ({ shipment }) => {
  // Simple coordinate mapping based on fixed layout
  const positions = {
    "Warehouse": { top: '50%', left: '50%' },
    "Outlet-1": { top: '5%', left: '10%' }, // Top Left
    "Outlet-2": { top: '5%', left: '80%' }, // Top Right
    "Outlet-3": { top: '85%', left: '10%' }, // Bottom Left
    "Outlet-4": { top: '85%', left: '80%' }, // Bottom Right
  };

  const start = positions["Warehouse"];
  const end = positions[shipment.to] || positions["Outlet-1"];

  return (
    <motion.div
      initial={{ top: start.top, left: start.left, opacity: 1, scale: 0.5 }}
      animate={{ top: end.top, left: end.left, opacity: 0, scale: 1.5 }}
      transition={{ duration: 5, ease: "easeInOut" }}
      className="absolute z-50 flex items-center justify-center pointer-events-none"
    >
      <div className="bg-yellow-400 text-black text-xs font-bold px-3 py-1 rounded-full shadow-lg border-2 border-white">
        📦 {shipment.qty} {shipment.item}
      </div>
    </motion.div>
  );
};

const OutletCard = ({ outlet, id }) => {
  if (!outlet) return null;
  return (
    <div className="bg-slate-800 p-4 rounded-xl border border-slate-700 shadow-lg hover:border-blue-500 transition-colors relative">
      <div className="absolute -top-2 -right-2 bg-blue-600 text-xs px-2 py-1 rounded">ID: {id}</div>
      <h2 className="text-lg font-bold text-blue-300 mb-2">{outlet.Name}</h2>
      <div className="space-y-1 text-sm">
        {Object.entries(outlet.Inventory).map(([item, qty]) => (
          <div key={item} className="flex justify-between">
            <span className="text-slate-400">{item}</span>
            <span className={qty < 20 ? "text-red-400 font-bold" : "text-slate-200"}>{qty}</span>
          </div>
        ))}
      </div>
    </div>
  );
};

const WarehouseCard = ({ warehouse }) => {
  return (
    <div className="bg-slate-800 p-6 rounded-2xl border-2 border-emerald-500/50 shadow-2xl shadow-emerald-500/10 w-80 text-center relative">
      <div className="absolute -top-3 left-1/2 transform -translate-x-1/2 bg-emerald-600 px-3 py-1 rounded-full text-xs font-bold uppercase tracking-wider">
        AI HUB
      </div>
      <h2 className="text-2xl font-bold text-emerald-300 mb-4">{warehouse.Name}</h2>
      <div className="space-y-2 text-sm text-left">
        {Object.entries(warehouse.Inventory).map(([item, qty]) => (
          <div key={item} className="flex justify-between border-b border-white/5 pb-1">
            <span className="text-slate-400">{item}</span>
            <span className="text-emerald-200 font-mono">{qty}</span>
          </div>
        ))}
      </div>
    </div>
  );
};

export default Dashboard;
