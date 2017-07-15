import { Component, OnInit, ViewEncapsulation, Input } from '@angular/core';
import { ImHistoryMessage, SocketService } from '../../providers';

@Component({
  selector: 'app-im-history-message',
  templateUrl: './im-history-message.component.html',
  styleUrls: ['./im-history-message.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class ImHistoryMessageComponent implements OnInit {
  selfId = '';
  @Input() chat: ImHistoryMessage = null;
  constructor(private socket: SocketService) {
    this.selfId = this.socket.key;
  }

  ngOnInit() {
    console.log('history message:', this.chat);
  }

}
