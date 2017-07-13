export interface ImHistoryMessage {
  type: HistoryMessageType;
  msg: string;
}

export enum HistoryMessageType {
  MYMESSAGE,
  OTHERMESSAGE
}