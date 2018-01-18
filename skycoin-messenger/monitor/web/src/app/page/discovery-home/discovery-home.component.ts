import { Component, OnInit, ViewEncapsulation, ViewChild } from '@angular/core';
import { MatTooltip } from '@angular/material';
import { ApiService } from '../../service';

@Component({
  selector: 'app-discovery-home',
  templateUrl: 'discovery-home.component.html',
  styleUrls: ['./discovery-home.component.scss'],
  encapsulation: ViewEncapsulation.None
})

export class DiscoveryHomeComponent implements OnInit {
  discoveryPubKey = '';
  @ViewChild('copyTooltip') tooltip: MatTooltip;

  constructor(private api: ApiService) { }
  copy(result: boolean) {
    if (result) {
      this.tooltip.disabled = false;
      this.tooltip.message = 'copied!';
      this.tooltip.hideDelay = 500;
      this.tooltip.show();
      setTimeout(() => {
        this.tooltip.disabled = true;
      }, 500);
    }
  }
  refresh(ev: Event) {
    this.init();
  }
  init() {
    this.api.getServerInfo().subscribe(key => {
      console.log('resp:', key);
      if (key) {
        this.discoveryPubKey = key;
      }
    });
  }
  ngOnInit() {
    this.init();
  }
}
