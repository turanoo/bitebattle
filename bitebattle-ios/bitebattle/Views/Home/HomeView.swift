import SwiftUI

struct HomeView: View {
    var body: some View {
        AppBackground {
            VStack(spacing: 32) {
                Spacer()
                
                // Logo
                Image(systemName: "fork.knife.circle.fill")
                    .resizable()
                    .scaledToFit()
                    .frame(width: 80, height: 80)
                    .foregroundColor(AppColors.orange)
                    .shadow(radius: 10)

                VStack(spacing: 18) {
                    NavigationLink(destination: PollsView()) {
                        AppButtonLabel(
                            title: "Polls",
                            icon: nil,
                            background: AppColors.brown.opacity(0.9),
                            foreground: AppColors.orange,
                            isLoading: false,
                            isDisabled: false
                        )
                    }
                    NavigationLink(destination: HeadToHeadView()) {
                        AppButtonLabel(
                            title: "Head-to-Head",
                            icon: nil,
                            background: AppColors.orange.opacity(0.8),
                            foreground: AppColors.textOnPrimary,
                            isLoading: false,
                            isDisabled: false
                        )
                    }
                }
                .padding(.horizontal, 24)

                Spacer()
            }
            .padding()
            .navigationTitle("BiteBattle")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarTrailing) {
                    NavigationLink(destination: AccountView()) {
                        HStack(spacing: 6) {
                            Image(systemName: "person.crop.circle")
                                .foregroundColor(AppColors.orange)
                            Text("Account")
                                .foregroundColor(AppColors.orange)
                        }
                        .padding(.vertical, 6)
                        .padding(.horizontal, 10)
                        .background(AppColors.surface.opacity(0.7))
                        .cornerRadius(8)
                    }
                }
            }
        }
    }
}

// Add a placeholder HeadToHeadView for navigation
struct HeadToHeadView: View {
    var body: some View {
        Text("Head-to-Head Coming Soon!")
            .font(.title)
            .foregroundColor(AppColors.orange)
            .padding()
    }
}

struct AppButtonLabel: View {
    let title: String
    let icon: String?
    let background: Color
    let foreground: Color
    let isLoading: Bool
    let isDisabled: Bool

    var body: some View {
        HStack {
            if let icon = icon {
                Image(systemName: icon)
                    .foregroundColor(foreground)
            }
            if isLoading {
                ProgressView()
                    .progressViewStyle(CircularProgressViewStyle(tint: foreground))
            } else {
                Text(title)
                    .fontWeight(.semibold)
                    .foregroundColor(foreground)
            }
        }
        .frame(maxWidth: .infinity, minHeight: 44)
        .padding(.vertical, 2)
        .background(isDisabled ? Color.gray.opacity(0.5) : background)
        .cornerRadius(12)
        .shadow(color: background.opacity(0.15), radius: 2, x: 0, y: 1)
        .scaleEffect(isDisabled ? 1.0 : 0.98, anchor: .center)
        .animation(.easeInOut(duration: 0.15), value: isDisabled)
    }
}
