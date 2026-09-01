import { useEffect, useState } from "react";
import { createRoot } from "react-dom/client";
import { Dashboard } from "./Dashboard";
import { InventoryOperationsScreen } from "./InventoryOperationsScreen";
import { LoginScreen } from "./LoginScreen";
import { ProductsScreen } from "./ProductsScreen";
import { PurchaseOrdersScreen } from "./PurchaseOrdersScreen";
import { WarehouseLocationScreen } from "./WarehouseLocationScreen";
import { stockpilotAPI } from "./api";
import type { User } from "./types";
import "./styles.css";
import "./lotExpiry.css";
type SessionState={kind:"checking"}|{kind:"signed-out"}|{kind:"signed-in";user:User};type View="dashboard"|"products"|"inventory"|"orders"|"warehouses";
function App(){const[session,setSession]=useState<SessionState>({kind:"checking"}),[view,setView]=useState<View>("dashboard");useEffect(()=>{let active=true;void stockpilotAPI.me().then(user=>{if(active)setSession({kind:"signed-in",user})},()=>{if(active)setSession({kind:"signed-out"})});return()=>{active=false}},[]);async function login(email:string,password:string){const result=await stockpilotAPI.login(email,password);setView("dashboard");setSession({kind:"signed-in",user:result.user})}async function logout(){try{await stockpilotAPI.logout()}finally{setView("dashboard");setSession({kind:"signed-out"})}}function expireSession(){setView("dashboard");setSession({kind:"signed-out"})}if(session.kind==="checking")return <main className="auth-page" aria-busy="true" aria-live="polite"><section className="auth-card session-check"><span className="brand-mark" aria-hidden="true">SP</span><p className="eyebrow">Secure workspace</p><h1>Checking your session</h1><p className="muted">StockPilot is verifying your signed-in session.</p></section></main>;if(session.kind==="signed-out")return <LoginScreen onLogin={login}/>;if(view==="products")return <ProductsScreen user={session.user} onBack={()=>setView("dashboard")} onSessionExpired={expireSession}/>;if(view==="inventory")return <InventoryOperationsScreen user={session.user} onBack={()=>setView("dashboard")} onSessionExpired={expireSession}/>;if(view==="orders")return <PurchaseOrdersScreen user={session.user} onBack={()=>setView("dashboard")} onSessionExpired={expireSession}/>;if(view==="warehouses")return <WarehouseLocationScreen user={session.user} onBack={()=>setView("dashboard")} onSessionExpired={expireSession}/>;return <Dashboard user={session.user} onLogout={logout} onSessionExpired={expireSession} onOpenProducts={()=>setView("products")} onOpenInventory={()=>setView("inventory")} onOpenOrders={()=>setView("orders")} onOpenWarehouses={()=>setView("warehouses")}/>}
const root=document.getElementById("root");if(!root)throw new Error("StockPilot root element was not found.");createRoot(root).render(<App/>);
