import { Component, OnInit, ViewEncapsulation,Input } from '@angular/core';

@Component({
  selector: 'app-im-history-view',
  templateUrl: './im-history-view.component.html',
  styleUrls: ['./im-history-view.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class ImHistoryViewComponent implements OnInit {
  @Input() chatList:Array<ImHistoryMessage>;
  @Input() from = 'self';
  constructor() { }
  ngOnInit() {
  }

}

export interface ImHistoryMessage {
  type: HistoryMessageType;
  msg: string;
}

export enum HistoryMessageType {
  MYMESSAGE,
  OTHERMESSAGE
}

