import SwiftUI

struct AppRouter: View {
    @AppStorage("isLoggedIn") var isLoggedIn: Bool = false

    var body: some View {
        NavigationView {
            if isLoggedIn {
                HomeView()
            } else {
                LandingView()
            }
        }
    }
}
