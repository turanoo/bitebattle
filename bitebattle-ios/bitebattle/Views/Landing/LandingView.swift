import SwiftUI

struct LandingView: View {
    var body: some View {
        AppBackground {
            AppGradientBackground {
                VStack(spacing: 32) {
                    Spacer()
                    // Logo placeholder
                    Image(systemName: "fork.knife.circle.fill")
                        .resizable()
                        .scaledToFit()
                        .frame(width: 100, height: 100)
                        .foregroundColor(AppColors.textOnPrimary)
                        .shadow(radius: 10)

                    Text("Welcome to BiteBattle")
                        .font(.largeTitle)
                        .fontWeight(.bold)
                        .foregroundColor(AppColors.textOnPrimary)
                        .multilineTextAlignment(.center)
                        .shadow(radius: 4)

                    VStack(spacing: 16) {
                        NavigationLink(destination: RegisterView()) {
                            AppButton(title: "Register", icon: nil, background: AppColors.surface.opacity(0.9), foreground: AppColors.orange, isLoading: false, isDisabled: false) {}
                        }
                        NavigationLink(destination: LoginView()) {
                            AppButton(title: "Login", icon: nil, background: AppColors.orange.opacity(0.8), foreground: AppColors.textOnPrimary, isLoading: false, isDisabled: false) {}
                        }
                    }
                    .padding(.horizontal, 24)

                    Spacer()
                }
                .padding()
            }
            .navigationTitle("Welcome")
            .navigationBarTitleDisplayMode(.inline)
        }
    }
}
