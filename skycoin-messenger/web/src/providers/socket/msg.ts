export interface ImHistoryMessage {
  From: string;
  Msg: string;
}

export interface RecentItem {
  name: string;
  last: string;
  unRead?: number;
}
