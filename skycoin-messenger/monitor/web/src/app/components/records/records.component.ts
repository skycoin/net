import { Component, OnInit } from '@angular/core';
import { ApiService, HashSig } from '../../service';
import { MatDialog } from '@angular/material';

@Component({
  selector: 'app-records',
  templateUrl: 'records.component.html',
  styleUrls: ['./records.component.scss']
})

export class RecordsComponent implements OnInit {
  nodeAddr = '';
  nodeKey = '';
  orders = [];
  balance = '0.000000';
  constructor(private api: ApiService, private dialog: MatDialog) { }

  ngOnInit() {
    this.getBalance();
    this.sig();
  }
  openWithdraw(ev: Event, content: any) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    this.dialog.open(content);
  }
  getBalance() {
    const hash = JSON.stringify({
      pubkey: this.nodeKey,
      timestamp: new Date().getTime,
    });
    this.api.getSig(this.nodeAddr, hash).subscribe(s => {
      const data = new FormData();
      data.append('data', hash);
      data.append('sig', s.sig);
      this.api.getBalance(data).subscribe(resp => {
        this.balance = resp.total;
      });
    });
  }
  sig() {
    const hash: HashSig = {
      pubkey: this.nodeKey,
      timestamp: new Date().getTime(),
    };
    this.api.getSig(this.nodeAddr, JSON.stringify(hash)).subscribe(resp => {
      const data = new FormData();
      data.append('data', JSON.stringify(hash));
      data.append('sig', resp.sig);
      this.api.getNodeOrders(data).subscribe((orders: Array<any>) => {
        this.orders = orders;
      });
    });
  }
}

