import { Component, OnInit, ViewEncapsulation, Input, ViewChild } from '@angular/core';
import { ImHistoryMessage, SocketService } from '../../providers';
import { MdMenuTrigger } from '@angular/material';

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
  }
  rightClick(ev: Event) {
    // ev.preventDefault();
    // this.contextMenu.openMenu();
  }
}
