import { useState } from "react";
import { ProductsScreen } from "./ProductsScreen";
import { ProductImportPanel } from "./ProductImportPanel";
import type { User } from "./types";
import "./productImport.css";

type Props = { user: User; onBack: () => void; onSessionExpired: () => void };

export function ProductsWorkspaceScreen({ user, onBack, onSessionExpired }: Props) {
  const [importOpen, setImportOpen] = useState(false);
  const canWrite = user.role === "admin" || user.role === "manager";
  return <>
    <ProductsScreen user={user} onBack={onBack} onSessionExpired={onSessionExpired} />
    {canWrite && <button className="catalog-import-launcher" type="button" onClick={() => setImportOpen(true)} aria-label="Open product CSV import">Import CSV</button>}
    {canWrite && importOpen && <ProductImportPanel onClose={() => setImportOpen(false)} onSessionExpired={onSessionExpired} onImported={async () => undefined} />}
  </>;
}
