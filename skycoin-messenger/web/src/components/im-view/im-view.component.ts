import { Component, OnInit, ViewEncapsulation, Input, OnChanges, SimpleChanges } from '@angular/core';
import { SocketService } from '../../providers';
import { ImHistoryMessage, HistoryMessageType } from '../im-history-view/im-history-view.component';

@Component({
  selector: 'app-im-view',
  templateUrl: './im-view.component.html',
  styleUrls: ['./im-view.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class ImViewComponent implements OnInit, OnChanges {

  chatList: Array<ImHistoryMessage>;
  msg = '';
  @Input() chatting = '';
  constructor(private socket: SocketService) { }

  ngOnInit() {
    // this.chattingUser = this.socket.chattingUser;
  }

  ngOnChanges(changes: SimpleChanges) {
    for (let propName in changes) {
      let chng = changes[propName];
      let data = SocketService.chatHistorys.get(chng.currentValue);
      if (data) {
        this.chatList = data;
      }else {
        this.chatList = [];
      }
    }
  }
  send(ev: KeyboardEvent) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.returnValue = false;
    this.msg = this.msg.trim();
    if (this.msg.length > 0) {
      this.addChat();
    }
  }

  addChat() {
    if (this.chatList.length > 0) {
      this.chatList.unshift({ type: HistoryMessageType.MYMESSAGE, msg: this.msg });
    } else {
      this.chatList = this.chatList.concat({ type: HistoryMessageType.MYMESSAGE, msg: this.msg });
    }
    SocketService.chatHistorys.set(this.chatting, this.chatList);
    console.log('log map:', SocketService.chatHistorys);
    // this.socket.msg(this.msg);
    this.msg = '';
  }
}
