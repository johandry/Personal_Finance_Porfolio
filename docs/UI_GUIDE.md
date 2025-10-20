# Screenshots & UI Guide

## 📸 Application Preview

### Dashboard Page (`index.html`)

```text
┌─────────────────────────────────────────────────────────────┐
│  💰 Finance Portfolio                Dashboard | Manage     │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│              Financial Dashboard                            │
│     Track your net worth and monitor your financial health  │
│                                                             │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐     │
│  │ 📈       │  │ 📉       │  │ 💎       │  │ 📊       │    │
│  │ Total    │  │ Total    │  │ Net      │  │ Total    │     │
│  │ Assets   │  │ Debts    │  │ Worth    │  │ Profit/  │     │
│  │          │  │          │  │          │  │ Loss     │     │
│  │ $50,000  │  │ $15,000  │  │ $35,000  │  │ +$5,000  │     │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘     │
│                                                             │
│  ┌────────────────────────┐  ┌────────────────────────┐     │
│  │ Asset Distribution     │  │ Net Worth Overview     │     │
│  │                        │  │                        │     │
│  │     [Pie Chart]        │  │     [Bar Chart]        │     │
│  │                        │  │                        │     │
│  │  Stock: 40%            │  │  Invested vs Current   │     │
│  │  Property: 35%         │  │  vs Profit/Loss        │     │
│  │  Cash: 25%             │  │                        │     │
│  └────────────────────────┘  └────────────────────────┘     │
│                                                             │
│  ┌───────────────────────────────────────────────────┐      │
│  │ Recent Assets                       View All →    │      │
│  ├──────┬──────┬─────────┬────────┬───────┬──────────┤      │
│  │ Name │ Type │ Quantity│ Current│ Total │ Profit/L │      │
│  ├──────┼──────┼─────────┼────────┼───────┼──────────┤      │
│  │ AAPL │Stock │   10    │ $175.50│$1,755 │  +$255   │      │
│  │ Tesla│Stock │    5    │ $250.00│$1,250 │  +$500   │      │
│  │ House│Prop. │    1    │$300,000│  ...  │ +$50,000 │      │
│  └──────┴──────┴─────────┴────────┴───────┴──────────┘      │
│                                                             │
│  ┌────────────────────────────────────────────────────┐     │
│  │ Recent Debts                         View All →    │     │
│  ├──────────┬──────────┬──────────┬──────────┬────────┤     │
│  │   Name   │   Type   │Principal │  Current │ Rate   │     │
│  ├──────────┼──────────┼──────────┼──────────┼────────┤     │
│  │ Car Loan │   Loan   │ $25,000  │ $20,000  │  4.5%  │     │
│  │ Mortgage │ Mortgage │$200,000  │$180,000  │  3.2%  │     │
│  └──────────┴──────────┴──────────┴──────────┴────────┘     │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Manage Page (`manage.html`)

```text
┌─────────────────────────────────────────────────────────────┐
│  💰 Finance Portfolio                Dashboard | Manage     │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│              Manage Your Finances                           │
│      Create, update, and delete your assets and debts       │
│                                                             │
│  ┌─────────┬─────────┐                                      │
│  │ Assets  │  Debts  │                 [+ Add New Asset]    │
│  └─────────┴─────────┘                                      │
│                                                             │
│  ┌────────────────────────────────────────────────────────┐ │
│  │ Assets Table                                           │ │
│  ├──────┬──────┬─────┬────────┬────┬──────┬─────┬────┬────┤ │
│  │ Name │ Type │ Buy │Current │Qty │Total │ P/L │Date│Act.│ │
│  ├──────┼──────┼─────┼────────┼────┼──────┼─────┼────┼────┤ │
│  │ AAPL │Stock │$150 │ $175.50│ 10 │$1,755│+$255│Jan │E D │ │
│  │Tesla │Stock │$200 │ $250.00│  5 │$1,250│+$250│Feb │E D │ │
│  │House │Prop. │$250k│ $300k  │  1 │$300k │+$50k│2020│E D │ │
│  └──────┴──────┴─────┴────────┴────┴──────┴─────┴────┴────┘ │
│                                                             │
│                         E = Edit, D = Delete                │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Add/Edit Asset Modal

```text
┌────────────────────────────────────────┐
│  Add New Asset                      ✕  │
├────────────────────────────────────────┤
│                                        │
│  Name *              Type *            │
│  [Apple Stock___]    [Stock ▼]         │
│                                        │
│  Buy Price *         Current Value     │
│  [150.00_______]     [175.50_______]   │
│                                        │
│  Quantity *          Currency          │
│  [10___________]     [USD ▼]           │
│                                        │
│  Purchase Date *     Source            │
│  [2024-01-01___]     [Manual ▼]        │
│                                        │
│  ┌────────────────────────────────┐    │
│  │   [Cancel]    [Save Asset]     │    │
│  └────────────────────────────────┘    │
└────────────────────────────────────────┘
```

## 🎨 Color Scheme

### Summary Cards

- **Assets Card**: Blue theme with 📈 icon
- **Debts Card**: Gray theme with 📉 icon  
- **Net Worth Card**: Purple gradient with 💎 icon (highlighted)
- **Profit/Loss Card**: Green/Red with 📊 icon

### Asset Type Badges

- **Stock**: Light blue background, dark blue text
- **Property**: Light yellow background, brown text
- **Car**: Light indigo background, dark indigo text
- **Cash**: Light green background, dark green text
- **Investment**: Light pink background, dark pink text

### Debt Type Badges

- **Credit Card**: Light red background, dark red text
- **Loan**: Light yellow background, brown text
- **Mortgage**: Light indigo background, dark indigo text
- **Other**: Light gray background, dark gray text

## 📱 Responsive Design

### Desktop (> 768px)

- 4 summary cards in a row
- 2 charts side by side
- Full table view
- Modal: 600px width

### Tablet (481px - 768px)

- 2 summary cards per row
- 1 chart per row
- Scrollable tables
- Modal: 90% width

### Mobile (< 480px)

- 1 summary card per row
- 1 chart per row
- Horizontal scroll for tables
- Full-width modal

## 🎯 Interactive Elements

### Buttons

- **Primary**: Indigo background, white text
  - "+ Add New Asset"
  - "Save Asset"
  
- **Secondary**: Gray background, white text
  - "Cancel"
  
- **Edit**: Indigo, small size
  - Pencil icon or "Edit" text
  
- **Delete**: Red, small size
  - Trash icon or "Delete" text

### Tables

- **Header**: Light gray background
- **Rows**: White background, hover effect
- **Empty State**: Centered message
- **Loading State**: Centered spinner/text

### Charts

- **Pie Chart**:
  - Colorful segments
  - Legend at bottom
  - Hover shows value + percentage
  
- **Bar Chart**:
  - Three bars (Invested, Current, P/L)
  - Different colors
  - Hover shows exact value

### Notifications (Toast)

- **Success**: Green background, slides in from right
- **Error**: Red background, slides in from right
- **Auto-dismiss**: 3 seconds
- **Position**: Bottom right

## 🔄 User Flow

### Adding an Asset

1. Click "Manage Assets & Debts" in nav
2. Ensure on "Assets" tab
3. Click "+ Add New Asset" button
4. Modal opens
5. Fill form fields
6. Click "Save Asset"
7. Modal closes
8. Table refreshes
9. Toast notification appears
10. Asset visible in table

### Viewing Dashboard

1. Land on index.html
2. Summary cards load
3. Charts render
4. Recent tables populate
5. Auto-refresh every 30s

### Editing an Asset

1. Go to Manage page
2. Find asset in table
3. Click "Edit" button
4. Modal opens with data
5. Modify fields
6. Click "Save Asset"
7. Table updates
8. Toast confirms

### Deleting an Asset

1. Find asset in table
2. Click "Delete" button
3. Confirm dialog appears
4. Click "OK"
5. Item removed from table
6. Toast confirms deletion

## 💡 Visual Feedback

### Loading States

- Tables: "Loading..." text centered
- Buttons: Disabled state during submit
- Charts: Empty until data loads

### Empty States

- Assets: "No assets found. Add your first asset!"
- Debts: "No debts found."
- Charts: Don't render if no data

### Error States

- Toast notification with error message
- Form validation on required fields
- Red border on invalid inputs

### Success States

- Green toast notification
- Updated data in tables
- Charts refresh automatically

## 🎨 Typography

- **Headers**: Bold, larger size
- **Amounts**: Bold, prominent
- **Labels**: Medium weight, smaller
- **Tables**: Regular, readable
- **Buttons**: Medium weight

## 🌈 Animations

- **Modal**: Slide down fade in
- **Toast**: Slide from right
- **Cards**: Hover lift effect
- **Buttons**: Color transition on hover
- **Tables**: Row highlight on hover

---

This UI provides a clean, modern, and intuitive experience for managing personal finances!
