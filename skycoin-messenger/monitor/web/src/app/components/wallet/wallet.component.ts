import { Component, OnInit, ViewChild } from '@angular/core';
import { MatIconRegistry, MatTooltip } from '@angular/material';

@Component({
  selector: 'app-wallet',
  templateUrl: 'wallet.component.html',
  styleUrls: ['./wallet.component.scss']
})

export class WalletComponent implements OnInit {
  @ViewChild('copyTooltip') tooltip: MatTooltip;
  key = 'Unknown';
  balance = '0.00000';
  records = [];
  constructor(private icon: MatIconRegistry) {
  }

  ngOnInit() {
    this.icon.registerFontClassAlias('fa');
  }
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
}
