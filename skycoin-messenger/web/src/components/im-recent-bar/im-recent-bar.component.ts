import { Component, OnInit, ViewEncapsulation, Input, ViewChildren, QueryList, Output, EventEmitter } from '@angular/core';
import { ImRecentItemComponent } from '../im-recent-item/im-recent-item.component';
import { SocketService } from '../../providers';

@Component({
  selector: 'app-im-recent-bar',
  templateUrl: './im-recent-bar.component.html',
  styleUrls: ['./im-recent-bar.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class ImRecentBarComponent implements OnInit {
  chatting = '';
  @ViewChildren(ImRecentItemComponent) items: QueryList<ImRecentItemComponent>;
  @Input() list = [];
  constructor(private socket: SocketService) { }
  ngOnInit() {
  }

  selectItem(item: ImRecentItemComponent) {
    // this.chatting.emit(item);
    this.chatting = item.name;
    this.socket.chattingUser = item.name;
    const tmp = this.items.filter((el) => {
      return el.name !== item.name;
    });
    tmp.forEach(el => {
      el.active = false;
    });
  }
}
