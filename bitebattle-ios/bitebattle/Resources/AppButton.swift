// AppButton.swift
// Minimalistic, consistent button for BiteBattle
import SwiftUI



public struct AppButton: View {
    let title: String
    let icon: String?
    let background: Color
    let foreground: Color
    let isLoading: Bool
    let isDisabled: Bool
    let action: () -> Void

    public init(title: String, icon: String? = nil, background: Color = AppColors.primary, foreground: Color = AppColors.textOnPrimary, isLoading: Bool = false, isDisabled: Bool = false, action: @escaping () -> Void) {
        self.title = title
        self.icon = icon
        self.background = background
        self.foreground = foreground
        self.isLoading = isLoading
        self.isDisabled = isDisabled
        self.action = action
    }

    public var body: some View {
        Button(action: action) {
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
        .disabled(isDisabled || isLoading)
    }
}
