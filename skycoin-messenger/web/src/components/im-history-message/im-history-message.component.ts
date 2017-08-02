import { Component, ViewEncapsulation, Input, ViewChild, OnInit } from '@angular/core';
import { ImHistoryMessage, SocketService } from '../../providers';
import { MdMenuTrigger } from '@angular/material';
import * as emojione from 'emojione';

@Component({
  selector: 'app-im-history-message',
  templateUrl: './im-history-message.component.html',
  styleUrls: ['./im-history-message.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class ImHistoryMessageComponent implements OnInit {
  selfId = '';
  @ViewChild(MdMenuTrigger) contextMenu: MdMenuTrigger;
  @Input() index: number;
  @Input() chat: ImHistoryMessage = null;
  constructor(private socket: SocketService) {
    this.selfId = this.socket.key;
  }
  ngOnInit() {
    this.chat.Msg = emojione.shortnameToImage(this.chat.Msg);
  }
  rightClick(ev: Event) {
    // ev.preventDefault();
    // this.contextMenu.openMenu();
  }
}
