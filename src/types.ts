export type Ship = {
  ship_id: number;
  name: string;
  description?: string | null;
  capacity?: number | null;
  length?: number | null;
  width?: number | null;
  draft?: number | null;
  cranes?: number | null;
  containers?: number | null;
  photo_url?: string | null;
}

export type RequestShip = {
  request_ship_id: number;
  status: string;
  comment?: string | null;
  containers_20ft_count: number;
  containers_40ft_count: number;
  ships?: Array<{ ship_id: number; ships_count: number; ship: Ship }>;
  loading_time?: number | null;
}
