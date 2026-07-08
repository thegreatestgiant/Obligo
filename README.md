# Obligo: Purpose-Driven Giving

For thousands of years, the biblical principle of setting aside a dedicated portion of our income for charity has guided people toward a life of generosity and stewardship. Whether you view it as a traditional tithe or a personal commitment to doing good, giving a percentage of what you earn ensures that charity remains a priority, not an afterthought.

But in today's world of fluctuating incomes, direct deposits, and digital donations, tracking exactly what you’ve earned versus what you’ve given can quickly become a messy spreadsheet of guesswork. 

**Obligo** is built to solve that. 

Obligo is a private, self-hosted web application designed to help you track your income and charitable giving side-by-side. It provides a persistent, live answer to a single important question: *"Am I keeping up with my giving goals?"*

### 🌐 Try the Live Demo
Experience the app yourself here: **[https://obligo.mk0.pp.ua/](https://obligo.mk0.pp.ua/)**

---

## What It Does

- **Eliminates the Guesswork:** Log paychecks as they come in and donations as you make them. Obligo automatically handles the math, showing you exactly where you stand.
- **Customizable Commitments:** By default, Obligo assumes a standard 10% giving target (a traditional tithe), but this is fully customizable in your settings to match your personal convictions.
- **Live Progress Tracking:** Your dashboard updates in real-time. It displays your total earnings, your total given, and exactly how much is left to fulfill your personal giving obligation.
- **Historical Ledger:** Keep a clear, paginated history of every contribution and paycheck you've logged, giving you complete peace of mind at the end of the year.
- **Complete Privacy:** Your finances are deeply personal. Obligo is designed for a single household or individual. Because it is self-hosted, your data remains entirely on your own infrastructure—no corporate tracking, no shared servers.

## How to Use It

1. **Create an account** and log in.
2. **Set your giving target** in the Settings (or keep the 10% default).
3. **Log your entries.** Add a "paycheck" whenever you earn money, and a "donation" whenever you give. 
4. **Watch your dashboard.** Obligo will dynamically calculate your fulfilled commitments and show you your running balances.

*(Note on historical accuracy: Your obligation is calculated based on your target percentage at the exact moment you log a paycheck. If you decide to increase your giving percentage later, it only applies to future earnings, ensuring your past history remains accurate.)*

---

## For Developers & System Admins

Obligo is open-source and meant to be self-hosted on infrastructure you control. It ships as a lightweight, single Docker image with database provisioning handled via Docker Compose.

* **Want to install and run Obligo?** Read the [Deployment Guide (DEPLOY.md)](./DEPLOY.md).
* **Want to contribute to the code or understand the architecture?** Read the [Developer Notes (DEVELOPERS.md)](./DEVELOPERS.md).
