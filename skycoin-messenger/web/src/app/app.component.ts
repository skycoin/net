import { Component, OnInit, ViewChild } from '@angular/core';
import { SocketService, ImHistoryMessage } from '../providers';
import { ImRecentItemComponent } from '../components';
@Component({
  selector: 'app-im',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  chatWindow = false;
  recent_list = []
  constructor(private socket: SocketService) {
  }
  ngOnInit(): void {
    // TODO Test User
    const testUser = window.location.search.slice(6);
    if (testUser !== '') {
      this.recent_list.push(testUser);
    }
    this.socket.chatHistorys.subscribe((data: Map<string, Array<ImHistoryMessage>>) => {
      data.forEach((value, key) => {
        if (this.recent_list.indexOf(key) <= -1) {
          this.recent_list.push(key.toLocaleUpperCase());
        }
      })
    })
    if (this.socket.chattingUser) {
      this.chatWindow = true;
    }
    this.socket.start();
  }

}
