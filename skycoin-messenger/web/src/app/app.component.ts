import { Component, OnInit, ViewChild } from '@angular/core';
import { SocketService, ImHistoryMessage, ToolService } from '../providers';
import { ImRecentItemComponent } from '../components';

@Component({
  selector: 'app-im',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  constructor(private socket: SocketService, private tool: ToolService) {
  }
  ngOnInit(): void {
    this.socket.start();
  }

}
