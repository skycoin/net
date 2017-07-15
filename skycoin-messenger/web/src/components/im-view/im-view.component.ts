import { Component, OnInit, ViewEncapsulation, Input, OnChanges, SimpleChanges } from '@angular/core';
import { SocketService } from '../../providers';
import { ImHistoryMessage } from '../../providers';

@Component({
  selector: 'app-im-view',
  templateUrl: './im-view.component.html',
  styleUrls: ['./im-view.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class ImViewComponent implements OnInit, OnChanges {
  chatList: Array<ImHistoryMessage>;
  historys: Map<string, Array<ImHistoryMessage>> = null;
  msg = '';
  @Input() chatting = '';
  constructor(private socket: SocketService) { }

  ngOnInit() {
    this.socket.chatHistorys.subscribe((data: Map<string, Array<ImHistoryMessage>>) => {
      console.log('send msg after:', data);
      this.historys = data;
    })
  }

  ngOnChanges(changes: SimpleChanges) {
    for (const propName in changes) {
      if (changes.hasOwnProperty(propName)) {
        const chng = changes[propName];
        if (!this.historys) {
          continue;
        }
        const data = this.historys.get((<string>chng.currentValue).toLocaleLowerCase());
        if (data) {
          this.chatList = data;
        } else {
          this.chatList = [];
        }
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
    this.socket.msg(this.chatting, this.msg);
    this.chatList.push({ From: this.socket.key, Msg: this.msg });
    this.socket.saveHistorys(this.chatting, this.socket.key, this.chatList);
    // this.chatList = this.historys.get(this.chatting);
    this.msg = '';
  }
}
